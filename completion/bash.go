package completion

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
)

// Bash is a bash completion generator.
type Bash struct{}

// Run generates bash completion script.
func (b Bash) Run(ctx *kong.Context) error {
	buf := new(bytes.Buffer)
	writePreamble(buf, ctx.Model.Name)
	b.gen(buf, ctx.Model.Node)
	writePostscript(buf, ctx.Model.Name)

	_, err := fmt.Fprint(ctx.Stdout, buf.String())
	if err != nil {
		return fmt.Errorf("unable to generate bash completion: %v", err)
	}
	return nil
}

// ShellCompDirective is a bit map representing the different behaviors the shell
// can be instructed to have once completions have been provided.
type ShellCompDirective int

const (
	// ShellCompDirectiveError indicates an error occurred and completions should be ignored.
	ShellCompDirectiveError ShellCompDirective = 1 << iota

	// ShellCompDirectiveNoSpace indicates that the shell should not add a space
	// after the completion even if there is a single completion provided.
	ShellCompDirectiveNoSpace

	// ShellCompDirectiveNoFileComp indicates that the shell should not provide
	// file completion even when no completion is provided.
	ShellCompDirectiveNoFileComp

	// ShellCompDirectiveFilterFileExt indicates that the provided completions
	// should be used as file extension filters.
	// For flags, using Command.MarkFlagFilename() and Command.MarkPersistentFlagFilename()
	// is a shortcut to using this directive explicitly.  The BashCompFilenameExt
	// annotation can also be used to obtain the same behavior for flags.
	ShellCompDirectiveFilterFileExt

	// ShellCompDirectiveFilterDirs indicates that only directory names should
	// be provided in file completion.  To request directory names within another
	// directory, the returned completions should specify the directory within
	// which to search.  The BashCompSubdirsInDir annotation can be used to
	// obtain the same behavior but only for flags.
	ShellCompDirectiveFilterDirs

	// ===========================================================================
	//
	// All directives using iota should be above this one.
	// For internal use.
	shellCompDirectiveMaxValue //nolint:deadcode,unused,varcheck

	// ShellCompDirectiveDefault indicates to let the shell perform its default
	// behavior after completions have been provided.
	// This one must be last to avoid messing up the iota count.
	ShellCompDirectiveDefault ShellCompDirective = 0
)

// Annotations for Bash completion.
const (
	// ShellCompNoDescRequestCmd is the name of the hidden command that is used to request
	// completion results without their description.  It is used by the shell completion scripts.
	ShellCompNoDescRequestCmd = "completion completeNoDesc"
	BashCompFilenameExt       = "kong_annotation_bash_completion_filename_extensions"
	BashCompCustom            = "kong_annotation_bash_completion_custom"
	BashCompOneRequiredFlag   = "kong_annotation_bash_completion_one_required_flag"
	BashCompSubdirsInDir      = "kong_annotation_bash_completion_subdirs_in_dir"

	activeHelpEnvVarSuffix = "_ACTIVE_HELP"
)

// activeHelpEnvVar returns the name of the program-specific ActiveHelp environment
// variable.  It has the format <PROGRAM>_ACTIVE_HELP where <PROGRAM> is the name of the
// root command in upper case, with all - replaced by _.
func activeHelpEnvVar(name string) string {
	// This format should not be changed: users will be using it explicitly.
	activeHelpEnvVar := strings.ToUpper(fmt.Sprintf("%s%s", name, activeHelpEnvVarSuffix))
	return strings.ReplaceAll(activeHelpEnvVar, "-", "_")
}

func writePreamble(buf io.StringWriter, name string) {
	writeString(buf, fmt.Sprintf("# bash completion for %-36s -*- shell-script -*-\n", name))
	writeString(buf, fmt.Sprintf(`
__%[1]s_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE:-} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__%[1]s_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__%[1]s_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__%[1]s_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__%[1]s_handle_go_custom_completion()
{
    __%[1]s_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=%[3]d
    local shellCompDirectiveNoSpace=%[4]d
    local shellCompDirectiveNoFileComp=%[5]d
    local shellCompDirectiveFilterFileExt=%[6]d
    local shellCompDirectiveFilterDirs=%[7]d

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly %[1]s allows to handle aliases
    args=("${words[@]:1}")
    # Disable ActiveHelp which is not supported for bash completion v1
    requestComp="%[8]s=0 ${words[0]} %[2]s ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __%[1]s_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __%[1]s_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __%[1]s_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __%[1]s_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __%[1]s_debug "${FUNCNAME[0]}: the completions are: ${out}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __%[1]s_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __%[1]s_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __%[1]s_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subdir
        # Use printf to strip any trailing newline
        subdir=$(printf "%%s" "${out}")
        if [ -n "$subdir" ]; then
            __%[1]s_debug "Listing directories in $subdir"
            __%[1]s_handle_subdirs_in_dir_flag "$subdir"
        else
            __%[1]s_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out}" -- "$cur")
    fi
}

__%[1]s_handle_reply()
{
    __%[1]s_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%%=*}"
                __%[1]s_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION:-}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi

            if [[ -z "${flag_parsing_disabled}" ]]; then
                # If flag parsing is enabled, we have completed the flags and can return.
                # If flag parsing is disabled, we may not know all (or any) of the flags, so we fallthrough
                # to possibly call handle_go_custom_completion.
                return 0;
            fi
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __%[1]s_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __%[1]s_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        if declare -F __%[1]s_custom_func >/dev/null; then
            # try command name qualified custom func
            __%[1]s_custom_func
        else
            # otherwise fall back to unqualified for compatibility
            declare -F __custom_func >/dev/null && __custom_func
        fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__%[1]s_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__%[1]s_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__%[1]s_handle_flag()
{
    __%[1]s_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue=""
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __%[1]s_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __%[1]s_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __%[1]s_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __%[1]s_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        __%[1]s_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__%[1]s_handle_noun()
{
    __%[1]s_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __%[1]s_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __%[1]s_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__%[1]s_handle_command()
{
    __%[1]s_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_%[1]s_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __%[1]s_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__%[1]s_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __%[1]s_handle_reply
        return
    fi
    __%[1]s_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __%[1]s_handle_flag
    elif __%[1]s_contains_word "${words[c]}" "${commands[@]}"; then
        __%[1]s_handle_command
    elif [[ $c -eq 0 ]]; then
        __%[1]s_handle_command
    elif __%[1]s_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __%[1]s_handle_command
        else
            __%[1]s_handle_noun
        fi
    else
        __%[1]s_handle_noun
    fi
    __%[1]s_handle_word
}

`, name, ShellCompNoDescRequestCmd,
		ShellCompDirectiveError, ShellCompDirectiveNoSpace, ShellCompDirectiveNoFileComp,
		ShellCompDirectiveFilterFileExt, ShellCompDirectiveFilterDirs, activeHelpEnvVar(name)))
}

func writePostscript(buf io.StringWriter, name string) {
	name = strings.ReplaceAll(name, ":", "__")
	writeString(buf, fmt.Sprintf("__start_%s()\n", name))
	writeString(buf, fmt.Sprintf(`{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __%[1]s_init_completion -n "=" || return
    fi

    local c=0
    local flag_parsing_disabled=
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("%[1]s")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function=""
    local last_command=""
    local nouns=()
    local noun_aliases=()

    __%[1]s_handle_word
}

`, name))
	writeString(buf, fmt.Sprintf(`if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_%s %s
else
    complete -o default -o nospace -F __start_%s %s
fi

`, name, name, name, name))
	writeString(buf, "# ex: ts=4 sw=4 et filetype=sh\n")
}

func writeCommands(buf io.StringWriter, cmd *kong.Node) {
	writeString(buf, "    commands=()\n")
	for _, c := range cmd.Children {
		if c == nil || c.Hidden {
			continue
		}
		writeString(buf, fmt.Sprintf("    commands+=(%q)\n", c.Name))
		writeCmdAliases(buf, c)
	}
	writeString(buf, "\n")
}

func writeFlagHandler(buf io.StringWriter, name string, annotations map[string][]string, cmd *kong.Node) {
	for key, value := range annotations {
		switch key {
		case BashCompFilenameExt:
			writeString(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) > 0 {
				ext = fmt.Sprintf("__%s_handle_filename_extension_flag ", cmd.Parent.Name) + strings.Join(value, "|")
			} else {
				ext = "_filedir"
			}
			writeString(buf, fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		case BashCompCustom:
			writeString(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			if len(value) > 0 {
				handlers := strings.Join(value, "; ")
				writeString(buf, fmt.Sprintf("    flags_completion+=(%q)\n", handlers))
			} else {
				writeString(buf, "    flags_completion+=(:)\n")
			}
		case BashCompSubdirsInDir:
			writeString(buf, fmt.Sprintf("    flags_with_completion+=(%q)\n", name))

			var ext string
			if len(value) == 1 {
				ext = fmt.Sprintf("__%s_handle_subdirs_in_dir_flag ", cmd.Parent.Name) + value[0]
			} else {
				ext = "_filedir -d"
			}
			writeString(buf, fmt.Sprintf("    flags_completion+=(%q)\n", ext))
		}
	}
}

const cbn = "\")\n"

func writeShortFlag(buf io.StringWriter, flag *kong.Flag, cmd *kong.Node) {
	name := fmt.Sprintf("%c", flag.Short)
	format := "    "
	if len(flag.DefaultValue.String()) == 0 {
		format += "two_word_"
	}
	format += "flags+=(\"-%s" + cbn
	writeString(buf, fmt.Sprintf(format, name))
	writeFlagHandler(buf, "-"+name, map[string][]string{}, cmd)
}

func writeFlag(buf io.StringWriter, flag *kong.Flag, cmd *kong.Node) {
	name := flag.Name
	format := "    flags+=(\"--%s"
	if len(flag.DefaultValue.String()) == 0 {
		format += "="
	}
	format += cbn
	writeString(buf, fmt.Sprintf(format, name))
	if len(flag.DefaultValue.String()) == 0 {
		format = "    two_word_flags+=(\"--%s" + cbn
		writeString(buf, fmt.Sprintf(format, name))
	}
	writeFlagHandler(buf, "--"+name, map[string][]string{}, cmd)
}

//nolint:deadcode,unused
func writeLocalNonPersistentFlag(buf io.StringWriter, flag *kong.Flag) {
	name := flag.Name
	format := "    local_nonpersistent_flags+=(\"--%[1]s" + cbn
	if len(flag.DefaultValue.String()) == 0 {
		format += "    local_nonpersistent_flags+=(\"--%[1]s=" + cbn
	}
	writeString(buf, fmt.Sprintf(format, name))
	if flag.Short > 0 {
		writeString(buf, fmt.Sprintf("    local_nonpersistent_flags+=(\"-%c\")\n", flag.Short))
	}
}

func writeFlags(buf io.StringWriter, cmd *kong.Node) {
	writeString(buf, `    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

`)

	for _, flag := range cmd.Flags {
		if nonCompletableFlag(flag) {
			continue
		}
		writeFlag(buf, flag, cmd)
		if flag.Short != 0 {
			writeShortFlag(buf, flag, cmd)
		}
	}

	writeString(buf, "\n")
}

func writeCmdAliases(buf io.StringWriter, cmd *kong.Node) {
	if len(cmd.Aliases) == 0 {
		return
	}

	sort.Strings(cmd.Aliases)

	writeString(buf, fmt.Sprint(`    if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then`, "\n"))
	for _, value := range cmd.Aliases {
		writeString(buf, fmt.Sprintf("        command_aliases+=(%q)\n", value))
		writeString(buf, fmt.Sprintf("        aliashash[%q]=%q\n", value, cmd.Name))
	}
	writeString(buf, `    fi`)
	writeString(buf, "\n")
}
func writeArgAliases(buf io.StringWriter, cmd *kong.Node) {
	writeString(buf, "    noun_aliases=()\n")
	sort.Strings(cmd.Aliases)
	for _, value := range cmd.Aliases {
		writeString(buf, fmt.Sprintf("    noun_aliases+=(%q)\n", value))
	}
}

func (b Bash) gen(buf io.StringWriter, cmd *kong.Node) {
	for _, c := range cmd.Children {
		if c == nil || c.Hidden {
			continue
		}
		b.gen(buf, c)
	}
	commandName := cmd.FullPath()
	commandName = strings.ReplaceAll(commandName, " ", "_")
	commandName = strings.ReplaceAll(commandName, ":", "__")

	if cmd.Parent == nil {
		writeString(buf, fmt.Sprintf("_%s_root_command()\n{\n", commandName))
	} else {
		writeString(buf, fmt.Sprintf("_%s()\n{\n", commandName))
	}

	writeString(buf, fmt.Sprintf("    last_command=%q\n", commandName))
	writeString(buf, "\n")
	writeString(buf, "    command_aliases=()\n")
	writeString(buf, "\n")

	writeCommands(buf, cmd)
	writeFlags(buf, cmd)
	writeArgAliases(buf, cmd)
	writeString(buf, "}\n\n")
}
