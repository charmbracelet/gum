package completion

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
)

type Complete struct {
	Arg []string `arg:"" passthrough:""`
}

type flagCompError struct {
	subCommand string
	flagName   string
}

func (e *flagCompError) Error() string {
	return "Subcommand '" + e.subCommand + "' does not support flag '" + e.flagName + "'"
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

	// All directives using iota should be above this one.
	// For internal use.
	shellCompDirectiveMaxValue

	// ShellCompDirectiveDefault indicates to let the shell perform its default
	// behavior after completions have been provided.
	// This one must be last to avoid messing up the iota count.
	ShellCompDirectiveDefault ShellCompDirective = 0
)

// Annotations for Bash completion.
const (
	// ShellCompRequestCmd is the name of the hidden command that is used to request
	// completion results from the program.  It is used by the shell completion scripts.
	ShellCompRequestCmd = "completion complete"
	// ShellCompNoDescRequestCmd is the name of the hidden command that is used to request
	// completion results without their description.  It is used by the shell completion scripts.
	ShellCompNoDescRequestCmd = "completion completeNoDesc"
	BashCompFilenameExt       = "kong_annotation_bash_completion_filename_extensions"
	BashCompCustom            = "kong_annotation_bash_completion_custom"
	BashCompOneRequiredFlag   = "kong_annotation_bash_completion_one_required_flag"
	BashCompSubdirsInDir      = "kong_annotation_bash_completion_subdirs_in_dir"
)

// Global map of flag completion functions. Make sure to use flagCompletionMutex before you try to read and write from it.
var flagCompletionFunctions = map[*kong.Flag]func(cmd *kong.Node, args []string, toComplete string) ([]string, ShellCompDirective){}

// lock for reading and writing from flagCompletionFunctions
var flagCompletionMutex = &sync.RWMutex{}

const (
	activeHelpMarker = "_activeHelp_ "
	// The below values should not be changed: programs will be using them explicitly
	// in their user documentation, and users will be using them explicitly.
	activeHelpEnvVarSuffix  = "_ACTIVE_HELP"
	activeHelpGlobalEnvVar  = "KONG_ACTIVE_HELP"
	activeHelpGlobalDisable = "0"
)

// activeHelpEnvVar returns the name of the program-specific ActiveHelp environment
// variable.  It has the format <PROGRAM>_ACTIVE_HELP where <PROGRAM> is the name of the
// root command in upper case, with all - replaced by _.
func activeHelpEnvVar(name string) string {
	// This format should not be changed: users will be using it explicitly.
	activeHelpEnvVar := strings.ToUpper(fmt.Sprintf("%s%s", name, activeHelpEnvVarSuffix))
	return strings.ReplaceAll(activeHelpEnvVar, "-", "_")
}

// WriteStringAndCheck writes a string into a buffer, and checks if the error is not nil.
func WriteStringAndCheck(b io.StringWriter, s string) {
	_, err := b.WriteString(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// Returns a string listing the different directive enabled in the specified parameter
func (d ShellCompDirective) string() string {
	var directives []string
	if d&ShellCompDirectiveError != 0 {
		directives = append(directives, "ShellCompDirectiveError")
	}
	if d&ShellCompDirectiveNoSpace != 0 {
		directives = append(directives, "ShellCompDirectiveNoSpace")
	}
	if d&ShellCompDirectiveNoFileComp != 0 {
		directives = append(directives, "ShellCompDirectiveNoFileComp")
	}
	if d&ShellCompDirectiveFilterFileExt != 0 {
		directives = append(directives, "ShellCompDirectiveFilterFileExt")
	}
	if d&ShellCompDirectiveFilterDirs != 0 {
		directives = append(directives, "ShellCompDirectiveFilterDirs")
	}
	if len(directives) == 0 {
		directives = append(directives, "ShellCompDirectiveDefault")
	}

	if d >= shellCompDirectiveMaxValue {
		return fmt.Sprintf("ERROR: unexpected ShellCompDirective value: %d", d)
	}
	return strings.Join(directives, ", ")
}

func (l Complete) BeforeApply(ctx *kong.Context) error {
	_, completions, directive, err := getCompletions(ctx.Model.Node, ctx.Args[2:])
	if err != nil {
		CompErrorln(err.Error())
		// Keep going for multiple reasons:
		// 1- There could be some valid completions even though there was an error
		// 2- Even without completions, we need to print the directive
	}

	noDescriptions := false // (cmd.CalledAs() == ShellCompNoDescRequestCmd)
	for _, comp := range completions {
		// if GetActiveHelpConfig(finalCmd) == activeHelpGlobalDisable {
		// 	// Remove all activeHelp entries in this case
		// 	if strings.HasPrefix(comp, activeHelpMarker) {
		// 		continue
		// 	}
		// }
		if noDescriptions {
			// Remove any description that may be included following a tab character.
			comp = strings.Split(comp, "\t")[0]
		}

		// Make sure we only write the first line to the output.
		// This is needed if a description contains a linebreak.
		// Otherwise the shell scripts will interpret the other lines as new flags
		// and could therefore provide a wrong completion.
		comp = strings.Split(comp, "\n")[0]

		// Finally trim the completion.  This is especially important to get rid
		// of a trailing tab when there are no description following it.
		// For example, a sub-command without a description should not be completed
		// with a tab at the end (or else zsh will show a -- following it
		// although there is no description).
		comp = strings.TrimSpace(comp)

		// Print each possible completion to stdout for the completion script to consume.
		fmt.Fprintln(ctx.Stdout, comp)
	}

	// As the last printout, print the completion directive for the completion script to parse.
	// The directive integer must be that last character following a single colon (:).
	// The completion script expects :<directive>
	fmt.Fprintf(ctx.Stdout, ":%d\n", directive)

	// Print some helpful info to stderr for the user to understand.
	// Output from stderr must be ignored by the completion script.
	fmt.Fprintf(ctx.Stderr, "Completion ended with directive: %s\n", directive.string())

	return nil
}

func parseFlags(c *kong.Node, args []string) (int, error) {
	flags := make([]*kong.Flag, 0)
	for _, f := range c.Flags {
		if f.Hidden || f.Tag.Optional || f.Tag.Ignored {
			continue
		}
		flags = append(flags, f)
	}
	return len(flags), nil
}

func traverse(c *kong.Node, args []string) (*kong.Node, []string, error) {
	if len(args) == 0 {
		return c, args, nil
	}

	// Find the sub-command to complete.
	for _, c := range c.Children {
		if c.Name == args[0] {
			return traverse(c, args[1:])
		}
	}

	// If we didn't find a sub-command, we are at the end of the path.
	// We can complete the sub-command name.
	return c, args, nil
}

func getCompletions(c *kong.Node, args []string) (*kong.Node, []string, ShellCompDirective, error) {
	// The last argument, which is not completely typed by the user,
	// should not be part of the list of arguments
	CompDebugln(c.Name, false)
	CompDebugln(fmt.Sprint(args), false)
	toComplete := args[len(args)-1]
	trimmedArgs := args[:len(args)-1]
	CompDebugln(fmt.Sprintf("toComplete: %v", toComplete), false)
	CompDebugln(fmt.Sprintf("trimmedArgs: %v", trimmedArgs), false)

	var finalCmd *kong.Node
	var finalArgs []string
	var err error

	finalCmd, finalArgs, err = traverse(c, trimmedArgs)
	if err != nil {
		// Unable to find the real command. E.g., <program> someInvalidCmd <TAB>
		return c, []string{}, ShellCompDirectiveDefault, fmt.Errorf("Unable to find a command for arguments: %v", trimmedArgs)
	}
	CompDebugln(fmt.Sprintf("finalCmd: %v", finalCmd.Name), false)
	CompDebugln(fmt.Sprintf("finalArgs: %v", finalArgs), false)
	// Check if we are doing flag value completion before parsing the flags.
	// This is important because if we are completing a flag value, we need to also
	// remove the flag name argument from the list of finalArgs or else the parsing
	// could fail due to an invalid value (incomplete) for the flag.
	flag, finalArgs, toComplete, flagErr := checkIfFlagCompletion(finalCmd, finalArgs, toComplete)

	// Check if interspersed is false or -- was set on a previous arg.
	// This works by counting the arguments. Normally -- is not counted as arg but
	// if -- was already set or interspersed is false and there is already one arg then
	// the extra added -- is counted as arg.
	flagCompletion := true
	newArgCount, _ := parseFlags(finalCmd, append(finalArgs, "--")) //len(finalCmd.Flags) //finalCmd.Flags().NArg()

	// Parse the flags early so we can check if required flags are set
	realArgCount, err := parseFlags(finalCmd, finalArgs)
	if err != nil {
		return finalCmd, []string{}, ShellCompDirectiveDefault, fmt.Errorf("Error while parsing flags from args %v: %s", finalArgs, err.Error())
	}

	if newArgCount > realArgCount {
		// don't do flag completion (see above)
		flagCompletion = false
	}
	// Error while attempting to parse flags
	if flagErr != nil {
		// If error type is flagCompError and we don't want flagCompletion we should ignore the error
		if _, ok := flagErr.(*flagCompError); !(ok && !flagCompletion) {
			return finalCmd, []string{}, ShellCompDirectiveDefault, flagErr
		}
	}

	// We only remove the flags from the arguments if DisableFlagParsing is not set.
	// This is important for commands which have requested to do their own flag completion.
	// if !finalCmd.DisableFlagParsing {
	// 	finalArgs = finalCmd.Flags().Args()
	// }

	if flag != nil && flagCompletion {
		// Check if we are completing a flag value subject to annotations
		// if validExts, present := flag.Annotations[BashCompFilenameExt]; present {
		// 	if len(validExts) != 0 {
		// 		// File completion filtered by extensions
		// 		return finalCmd, validExts, ShellCompDirectiveFilterFileExt, nil
		// 	}

		// The annotation requests simple file completion.  There is no reason to do
		// that since it is the default behavior anyway.  Let's ignore this annotation
		// in case the program also registered a completion function for this flag.
		// Even though it is a mistake on the program's side, let's be nice when we can.
		// }

		// if subDir, present := flag.Annotations[BashCompSubdirsInDir]; present {
		// 	if len(subDir) == 1 {
		// 		// Directory completion from within a directory
		// 		return finalCmd, subDir, ShellCompDirectiveFilterDirs, nil
		// 	}
		// 	// Directory completion
		// 	return finalCmd, []string{}, ShellCompDirectiveFilterDirs, nil
		// }
	}

	var completions []string
	var directive ShellCompDirective

	// Enforce flag groups before doing flag completions
	// finalCmd.enforceFlagGroupsForCompletion()

	// Note that we want to perform flagname completion even if finalCmd.DisableFlagParsing==true;
	// doing this allows for completion of persistent flag names even for commands that disable flag parsing.
	//
	// When doing completion of a flag name, as soon as an argument starts with
	// a '-' we know it is a flag.  We cannot use isFlagArg() here as it requires
	// the flag name to be complete
	if flag == nil && len(toComplete) > 0 && toComplete[0] == '-' && !strings.Contains(toComplete, "=") && flagCompletion {
		CompDebugln("Flag name completion", false)
		// First check for required flags
		completions = completeRequireFlags(finalCmd, toComplete)

		// If we have not found any required flags, only then can we show regular flags
		if len(completions) == 0 {
			for _, flag := range finalCmd.Flags {
				completions = append(completions, getFlagNameCompletions(flag, toComplete)...)
			}
		}

		directive = ShellCompDirectiveNoFileComp
		if len(completions) == 1 && strings.HasSuffix(completions[0], "=") {
			// If there is a single completion, the shell usually adds a space
			// after the completion.  We don't want that if the flag ends with an =
			directive = ShellCompDirectiveNoSpace
		}

		// if !finalCmd.DisableFlagParsing {
		// 	// If DisableFlagParsing==false, we have completed the flags as known by Cobra;
		// 	// we can return what we found.
		// 	// If DisableFlagParsing==true, Cobra may not be aware of all flags, so we
		// 	// let the logic continue to see if ValidArgsFunction needs to be called.
		// 	return finalCmd, completions, directive, nil
		// }
	} else {
		directive = ShellCompDirectiveDefault
		// if flag == nil {
		foundLocalNonPersistentFlag := false
		// If TraverseChildren is true on the root command we don't check for
		// local flags because we can use a local flag on a parent command
		// if !finalCmd.Root().TraverseChildren {
		if finalCmd.Parent != nil && len(finalCmd.Parent.Children) == 0 {
			// Check if there are any local, non-persistent flags on the command-line
			// localNonPersistentFlags := finalCmd.LocalNonPersistentFlags()
			// finalCmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
			// 	if localNonPersistentFlags.Lookup(flag.Name) != nil && flag.Changed {
			// 		foundLocalNonPersistentFlag = true
			// 	}
			// })
		}

		// Complete subcommand names, including the help command
		if len(finalArgs) == 0 && !foundLocalNonPersistentFlag {
			// We only complete sub-commands if:
			// - there are no arguments on the command-line and
			// - there are no local, non-persistent flags on the command-line or TraverseChildren is true
			for _, subCmd := range finalCmd.Children {
				// if subCmd.IsAvailableCommand() || subCmd == finalCmd.helpCommand {
				if !subCmd.Hidden {
					if strings.HasPrefix(subCmd.Name, toComplete) {
						completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name, subCmd.Help))
					}
					directive = ShellCompDirectiveNoFileComp
				}
			}
		}

		// Complete required flags even without the '-' prefix
		completions = append(completions, completeRequireFlags(finalCmd, toComplete)...)

		// Always complete ValidArgs, even if we are completing a subcommand name.
		// This is for commands that have both subcommands and ValidArgs.
		// if len(finalCmd.ValidArgs) > 0 {
		// 	if len(finalArgs) == 0 {
		// 		// ValidArgs are only for the first argument
		// 		for _, validArg := range finalCmd.ValidArgs {
		// 			if strings.HasPrefix(validArg, toComplete) {
		// 				completions = append(completions, validArg)
		// 			}
		// 		}
		// 		directive = ShellCompDirectiveNoFileComp

		// 		// If no completions were found within commands or ValidArgs,
		// 		// see if there are any ArgAliases that should be completed.
		// 		if len(completions) == 0 {
		// 			for _, argAlias := range finalCmd.ArgAliases {
		// 				if strings.HasPrefix(argAlias, toComplete) {
		// 					completions = append(completions, argAlias)
		// 				}
		// 			}
		// 		}
		// 	}

		// 	// If there are ValidArgs specified (even if they don't match), we stop completion.
		// 	// Only one of ValidArgs or ValidArgsFunction can be used for a single command.
		// 	return finalCmd, completions, directive, nil
		// }

		// Let the logic continue so as to add any ValidArgsFunction completions,
		// even if we already found sub-commands.
		// This is for commands that have subcommands but also specify a ValidArgsFunction.
		// }
	}

	// Find the completion function for the flag or command
	var completionFn func(cmd *kong.Node, args []string, toComplete string) ([]string, ShellCompDirective)
	if flag != nil && flagCompletion {
		flagCompletionMutex.RLock()
		completionFn = flagCompletionFunctions[flag]
		flagCompletionMutex.RUnlock()
	} else {
		completionFn = nil //finalCmd.ValidArgsFunction
	}
	if completionFn != nil {
		// Go custom completion defined for this flag or command.
		// Call the registered completion function to get the completions.
		var comps []string
		comps, directive = completionFn(finalCmd, finalArgs, toComplete)
		completions = append(completions, comps...)
	}

	return finalCmd, completions, directive, nil
}

func completeRequireFlags(finalCmd *kong.Node, toComplete string) []string {
	var completions []string

	// We cannot use finalCmd.Flags() because we may not have called ParsedFlags() for commands
	// that have set DisableFlagParsing; it is ParseFlags() that merges the inherited and
	// non-inherited flags.
	// finalCmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
	// 	doCompleteRequiredFlags(flag)
	// })
	// finalCmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
	// 	doCompleteRequiredFlags(flag)
	// })
	for _, f := range finalCmd.Flags {
		if f.Required {
			CompDebugln(fmt.Sprintf("flag: %s help: %s", f.Name, f.Help), false)
			completions = append(completions, getFlagNameCompletions(f, toComplete)...)
		}
	}

	return completions
}

func getFlagNameCompletions(flag *kong.Flag, toComplete string) []string {
	if nonCompletableFlag(flag) {
		return []string{}
	}

	var completions []string
	flagName := "--" + flag.Name
	if strings.HasPrefix(flagName, toComplete) {
		// Flag without the =
		flagHelp := flag.Help
		completions = append(completions, fmt.Sprintf("%s\t%s", flagName, flagHelp))

		// Why suggest both long forms: --flag and --flag= ?
		// This forces the user to *always* have to type either an = or a space after the flag name.
		// Let's be nice and avoid making users have to do that.
		// Since boolean flags and shortname flags don't show the = form, let's go that route and never show it.
		// The = form will still work, we just won't suggest it.
		// This also makes the list of suggested flags shorter as we avoid all the = forms.
		//
		// if len(flag.Default) == 0 {
		// 	// Flag requires a value, so it can be suffixed with =
		// 	flagName += "="
		// 	completions = append(completions, fmt.Sprintf("%s\t%s", flagName, flag.Help))
		// }
	}

	flagName = "-" + fmt.Sprintf("%c", flag.Short)
	if flag.Short != 0 && strings.HasPrefix(flagName, toComplete) {
		completions = append(completions, fmt.Sprintf("%s\t%s", flagName, flag.Help))
	}

	return completions
}

func isFlagArg(arg string) bool {
	return ((len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
}

func checkIfFlagCompletion(finalCmd *kong.Node, args []string, lastArg string) (*kong.Flag, []string, string, error) {
	var flagName string
	trimmedArgs := args
	flagWithEqual := false
	orgLastArg := lastArg

	// When doing completion of a flag name, as soon as an argument starts with
	// a '-' we know it is a flag.  We cannot use isFlagArg() here as that function
	// requires the flag name to be complete
	if len(lastArg) > 0 && lastArg[0] == '-' {
		if index := strings.Index(lastArg, "="); index >= 0 {
			// Flag with an =
			if strings.HasPrefix(lastArg[:index], "--") {
				// Flag has full name
				flagName = lastArg[2:index]
			} else {
				// Flag is shorthand
				// We have to get the last shorthand flag name
				// e.g. `-asd` => d to provide the correct completion
				// https://github.com/spf13/cobra/issues/1257
				flagName = lastArg[index-1 : index]
			}
			lastArg = lastArg[index+1:]
			flagWithEqual = true
		} else {
			// Normal flag completion
			return nil, args, lastArg, nil
		}
	}

	if len(flagName) == 0 {
		if len(args) > 0 {
			prevArg := args[len(args)-1]
			if isFlagArg(prevArg) {
				// Only consider the case where the flag does not contain an =.
				// If the flag contains an = it means it has already been fully processed,
				// so we don't need to deal with it here.
				if index := strings.Index(prevArg, "="); index < 0 {
					if strings.HasPrefix(prevArg, "--") {
						// Flag has full name
						flagName = prevArg[2:]
					} else {
						// Flag is shorthand
						// We have to get the last shorthand flag name
						// e.g. `-asd` => d to provide the correct completion
						// https://github.com/spf13/cobra/issues/1257
						flagName = prevArg[len(prevArg)-1:]
					}
					// Remove the uncompleted flag or else there could be an error created
					// for an invalid value for that flag
					trimmedArgs = args[:len(args)-1]
				}
			}
		}
	}

	if len(flagName) == 0 {
		// Not doing flag completion
		return nil, trimmedArgs, lastArg, nil
	}

	flag := findFlag(finalCmd, flagName)
	if flag == nil {
		// Flag not supported by this command, the interspersed option might be set so return the original args
		return nil, args, orgLastArg, &flagCompError{subCommand: finalCmd.Name, flagName: flagName}
	}

	if !flagWithEqual {
		if len(flag.Default) != 0 {
			// We had assumed dealing with a two-word flag but the flag is a boolean flag.
			// In that case, there is no value following it, so we are not really doing flag completion.
			// Reset everything to do noun completion.
			trimmedArgs = args
			flag = nil
		}
	}

	return flag, trimmedArgs, lastArg, nil
}

func shorthandLookup(flags []*kong.Flag, name string) *kong.Flag {
	for _, flag := range flags {
		if flag.Name == name {
			return flag
		}
	}
	return nil
}

func findFlag(cmd *kong.Node, name string) *kong.Flag {
	flagSet := cmd.Flags
	for _, flag := range flagSet {
		if flag.Name == name {
			return flag
		}
	}
	return nil
}

func nonCompletableFlag(flag *kong.Flag) bool {
	return flag.Hidden
}

// CompDebug prints the specified string to the same file as where the
// completion script prints its logs.
// Note that completion printouts should never be on stdout as they would
// be wrongly interpreted as actual completion choices by the completion script.
func CompDebug(msg string, printToStdErr bool) {
	msg = fmt.Sprintf("[Debug] %s", msg)

	// Such logs are only printed when the user has set the environment
	// variable BASH_COMP_DEBUG_FILE to the path of some file to be used.
	if path := os.Getenv("BASH_COMP_DEBUG_FILE"); path != "" {
		f, err := os.OpenFile(path,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			WriteStringAndCheck(f, msg)
		}
	}

	if printToStdErr {
		// Must print to stderr for this not to be read by the completion script.
		fmt.Fprint(os.Stderr, msg)
	}
}

// CompDebugln prints the specified string with a newline at the end
// to the same file as where the completion script prints its logs.
// Such logs are only printed when the user has set the environment
// variable BASH_COMP_DEBUG_FILE to the path of some file to be used.
func CompDebugln(msg string, printToStdErr bool) {
	CompDebug(fmt.Sprintf("%s\n", msg), printToStdErr)
}

// CompError prints the specified completion message to stderr.
func CompError(msg string) {
	msg = fmt.Sprintf("[Error] %s", msg)
	CompDebug(msg, true)
}

// CompErrorln prints the specified completion message to stderr with a newline at the end.
func CompErrorln(msg string) {
	CompError(fmt.Sprintf("%s\n", msg))
}
