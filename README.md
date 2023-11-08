Gum
===

<p>
    <a href="https://stuff.charm.sh/gum/nutritional-information.png" target="_blank"><img src="https://stuff.charm.sh/gum/gum.png" alt="Gum Image" width="450" /></a>
    <br><br>
    <a href="https://github.com/charmbracelet/gum/releases"><img src="https://img.shields.io/github/release/charmbracelet/gum.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/gum?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
    <a href="https://github.com/charmbracelet/gum/actions"><img src="https://github.com/charmbracelet/gum/workflows/build/badge.svg" alt="Build Status"></a>
</p>

A tool for glamorous shell scripts. Leverage the power of
[Bubbles](https://github.com/charmbracelet/bubbles) and [Lip
Gloss](https://github.com/charmbracelet/lipgloss) in your scripts and aliases
without writing any Go code!

<img alt="Shell running the ./demo.sh script" width="600" src="https://stuff.charm.sh/gum/demo.gif">

The above example is running from a single shell script ([source](./examples/demo.sh)).

## Tutorial

Gum provides highly configurable, ready-to-use utilities to help you write
useful shell scripts and dotfiles aliases with just a few lines of code.

Let's build a simple script to help you write [Conventional
Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) for your
dotfiles.

Start with a `#!/bin/sh`.
```bash
#!/bin/sh
```

Ask for the commit type with `gum choose`:

```bash
gum choose "fix" "feat" "docs" "style" "refactor" "test" "chore" "revert"
```

> Tip: this command itself will print to `stdout` which is not all that useful.
To make use of the command later on you can save the stdout to a `$VARIABLE` or
`file.txt`.

Prompt for an (optional) scope for the commit:

```bash
gum input --placeholder "scope"
```

Prompt for a commit message:

```bash
gum input --placeholder "Summary of this change"
```

Prompt for a detailed (multi-line) explanation of the changes:

```bash
gum write --placeholder "Details of this change (CTRL+D to finish)"
```

Prompt for a confirmation before committing:
> `gum confirm` exits with status `0` if confirmed and status `1` if cancelled.

```bash
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
```

Putting it all together...

```bash
#!/bin/sh
TYPE=$(gum choose "fix" "feat" "docs" "style" "refactor" "test" "chore" "revert")
SCOPE=$(gum input --placeholder "scope")

# Since the scope is optional, wrap it in parentheses if it has a value.
test -n "$SCOPE" && SCOPE="($SCOPE)"

# Pre-populate the input with the type(scope): so that the user may change it
SUMMARY=$(gum input --value "$TYPE$SCOPE: " --placeholder "Summary of this change")
DESCRIPTION=$(gum write --placeholder "Details of this change (CTRL+D to finish)")

# Commit these changes
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
```

<img alt="Running the ./examples/commit.sh script to commit to git" width="600" src="https://stuff.charm.sh/gum/commit_2.gif">

## Installation

Use a package manager:

```bash
# macOS or Linux
brew install gum

# Arch Linux (btw)
pacman -S gum

# Nix
nix-env -iA nixpkgs.gum
# Or, with flakes
nix run "github:charmbracelet/gum" -- --help

# Debian/Ubuntu
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install gum

# Fedora/RHEL
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
sudo yum install gum

# Alpine
apk add gum

# Android (via termux)
pkg install gum

# Windows (via WinGet or Scoop)
winget install charmbracelet.gum
scoop install charm-gum
```

Or download it:

* [Packages][releases] are available in Debian, RPM, and Alpine formats
* [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

Or just install it with `go`:

```bash
go install github.com/charmbracelet/gum@latest
```

[releases]: https://github.com/charmbracelet/gum/releases

## Customization

`gum` is designed to be embedded in scripts and supports all sorts of use
cases. Components are configurable and customizable to fit your theme and
use case.

You can customize with `--flags`. See `gum <command> --help` for a full view of
each command's customization and configuration options.

For example, let's use an `input` and change the cursor color, prompt color,
prompt indicator, placeholder text, width, and pre-populate the value:

```bash
gum input --cursor.foreground "#FF0" --prompt.foreground "#0FF" --prompt "* " \
    --placeholder "What's up?" --width 80 --value "Not much, hby?"
```

You can also use `ENVIRONMENT_VARIABLES` to customize `gum` by default, this is
useful to keep a consistent theme for all your `gum` commands.

```bash
export GUM_INPUT_CURSOR_FOREGROUND="#FF0"
export GUM_INPUT_PROMPT_FOREGROUND="#0FF"
export GUM_INPUT_PLACEHOLDER="What's up?"
export GUM_INPUT_PROMPT="* "
export GUM_INPUT_WIDTH=80

# Uses values configured through environment variables above but can still be
# overridden with flags.
gum input
```

<img alt="Gum input displaying most customization options" width="600" src="https://stuff.charm.sh/gum/customization.gif">

## Interaction

#### Input

Prompt for input with a simple command.

```bash
gum input > answer.txt
```

Prompt for sensitive input with the `--password` flag.

```bash
gum input --password > password.txt
```

<img src="https://stuff.charm.sh/gum/input_1.gif" width="600" alt="Shell running gum input typing Not much, you?" />

#### Write

Prompt for some multi-line text.

Note: `CTRL+D` is used to complete text entry. `CTRL+C` and `esc` will cancel.

```bash
gum write > story.txt
```

<img src="https://stuff.charm.sh/gum/write.gif" width="600" alt="Shell running gum write typing a story" />

#### Filter

Use fuzzy matching to filter a list of values:

```bash
echo Strawberry >> flavors.txt
echo Banana >> flavors.txt
echo Cherry >> flavors.txt
cat flavors.txt | gum filter > selection.txt
```

<img src="https://stuff.charm.sh/gum/filter.gif" width="600" alt="Shell running gum filter on different bubble gum flavors" />

You can also select multiple items with the `--limit` flag, which determines
the maximum number of items that can be chosen.

```bash
cat flavors.txt | gum filter --limit 2
```

Or, allow any number of selections with the `--no-limit` flag.

```bash
cat flavors.txt | gum filter --no-limit
```

#### Choose

Choose an option from a list of choices.

```bash
echo "Pick a card, any card..."
CARD=$(gum choose --height 15 {{A,K,Q,J},{10..2}}" "{♠,♥,♣,♦})
echo "Was your card the $CARD?"
```

You can also select multiple items with the `--limit` flag, which determines
the maximum of items that can be chosen.

```bash
echo "Pick your top 5 songs."
cat songs.txt | gum choose --limit 5
```

Or, allow any number of selections with the `--no-limit` flag.

```bash
echo "What do you need from the grocery store?"
cat foods.txt | gum choose --no-limit
```

<img src="https://stuff.charm.sh/gum/choose.gif" width="600" alt="Shell running gum choose with numbers and gum flavors" />

#### Confirm

Confirm whether to perform an action. Exits with code `0` (affirmative) or `1`
(negative) depending on selection.

```bash
gum confirm && rm file.txt || echo "File not removed"
```

<img src="https://stuff.charm.sh/gum/confirm_2.gif" width="600" alt="Shell running gum confirm" />

#### File

Prompt the user to select a file from the file tree.

```bash
EDITOR $(gum file $HOME)
```

<img src="https://stuff.charm.sh/gum/file.gif" width="600" alt="Shell running gum file" />

#### Pager

Scroll through a long document with line numbers and a fully customizable viewport.

```bash
gum pager < README.md
```

<img src="https://stuff.charm.sh/gum/pager.gif" width="600" alt="Shell running gum pager" />

#### Spin

Display a spinner while running a script or command. The spinner will
automatically stop after the given command exits.

To view or pipe the command's output, use the `--show-output` flag.

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

<img src="https://stuff.charm.sh/gum/spin.gif" width="600" alt="Shell running gum spin while sleeping for 5 seconds" />

Available spinner types include: `line`, `dot`, `minidot`, `jump`, `pulse`, `points`, `globe`, `moon`, `monkey`, `meter`, `hamburger`.

#### Table

Select a row from some tabular data.

```bash
gum table < flavors.csv | cut -d ',' -f 1
```

<img src="https://stuff.charm.sh/gum/table.gif" width="600" alt="Shell running gum table" />

## Styling

#### Style

Pretty print any string with any layout with one command.

```bash
gum style \
	--foreground 212 --border-foreground 212 --border double \
	--align center --width 50 --margin "1 2" --padding "2 4" \
	'Bubble Gum (1¢)' 'So sweet and so fresh!'
```

<img src="https://stuff.charm.sh/gum/style.gif" width="600" alt="Bubble Gum, So sweet and so fresh!" />

## Layout

#### Join

Combine text vertically or horizontally. Use this command with `gum style` to
build layouts and pretty output.

Tip: Always wrap the output of `gum style` in quotes to preserve newlines
(`\n`) when using it as an argument in the `join` command.

```bash
I=$(gum style --padding "1 5" --border double --border-foreground 212 "I")
LOVE=$(gum style --padding "1 4" --border double --border-foreground 57 "LOVE")
BUBBLE=$(gum style --padding "1 8" --border double --border-foreground 255 "Bubble")
GUM=$(gum style --padding "1 5" --border double --border-foreground 240 "Gum")

I_LOVE=$(gum join "$I" "$LOVE")
BUBBLE_GUM=$(gum join "$BUBBLE" "$GUM")
gum join --align center --vertical "$I_LOVE" "$BUBBLE_GUM"
```

<img src="https://stuff.charm.sh/gum/join.gif" width="600" alt="I LOVE Bubble Gum written out in four boxes with double borders around them." />

## Format

`format` processes and formats bodies of text. `gum format` can parse markdown,
template strings, and named emojis.

```bash
# Format some markdown
gum format -- "# Gum Formats" "- Markdown" "- Code" "- Template" "- Emoji"
echo "# Gum Formats\n- Markdown\n- Code\n- Template\n- Emoji" | gum format

# Syntax highlight some code
cat main.go | gum format -t code

# Render text any way you want with templates
echo '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' \
    | gum format -t template

# Display your favorite emojis!
echo 'I :heart: Bubble Gum :candy:' | gum format -t emoji
```

For more information on template helpers, see the [Termenv
docs](https://github.com/muesli/termenv#template-helpers). For a full list of
named emojis see the [GitHub API](https://api.github.com/emojis).

<img src="https://stuff.charm.sh/gum/format.gif" width="600" alt="Running gum format for different types of formats" />

## Log

`log` logs messages to the terminal at using different levels and styling using
the [`charmbracelet/log`](https://github.com/charmbracelet/log) library.

```bash
# Log some debug information.
gum log --structured --level debug "Creating file..." name file.txt
# DEBUG Unable to create file. name=temp.txt

# Log some error.
gum log --structured --level error "Unable to create file." name file.txt
# ERROR Unable to create file. name=temp.txt
```

See [`charmbracelet/log`](https://github.com/charmbracelet/log) for more usage.

<img src="https://vhs.charm.sh/vhs-6jupuFM0s2fXiUrBE0I1vU.gif" width="600" alt="Running gum log with debug and error levels" />

## Examples

See the [examples](./examples/) directory for more real world use cases.

How to use `gum` in your daily workflows:

#### Write a commit message

Prompt for input to write git commit messages with a short summary and
longer details with `gum input` and `gum write`.

Bonus points: use `gum filter` with the [Conventional Commits
Specification](https://www.conventionalcommits.org/en/v1.0.0/#summary) as a
prefix for your commit message.

```bash
git commit -m "$(gum input --width 50 --placeholder "Summary of changes")" \
           -m "$(gum write --width 80 --placeholder "Details of changes (CTRL+D to finish)")"
```

#### Open files in your `$EDITOR`

By default, `gum filter` will display a list of all files (searched
recursively) through your current directory, with some sensible ignore settings
(`.git`, `node_modules`). You can use this command to easily to pick a file and
open it in your `$EDITOR`.

```bash
$EDITOR $(gum filter)
```

#### Connect to a TMUX session

Pick from a running `tmux` session and attach to it. Or, if you're already in a
`tmux` session, switch sessions.

```bash
SESSION=$(tmux list-sessions -F \#S | gum filter --placeholder "Pick session...")
tmux switch-client -t $SESSION || tmux attach -t $SESSION
```

<img src="https://stuff.charm.sh/gum/pick-tmux-session.gif" width="600" alt="Picking a tmux session with gum filter" />

#### Pick commit hash from your Git history

Filter through your git history searching for commit messages, copying the
commit hash of the commit you select.

```bash
git log --oneline | gum filter | cut -d' ' -f1 # | copy
```

<img src="https://stuff.charm.sh/gum/pick-commit.gif" width="600" alt="Picking a commit with gum filter" />

#### Skate Passwords

Build a simple (encrypted) password selector with [Skate](https://github.com/charmbracelet/skate).

Save all your passwords to [Skate](https://github.com/charmbracelet/skate) with `skate set github@pass.db PASSWORD`, etc...

```
skate list -k | gum filter | xargs skate get
```

<img src="https://stuff.charm.sh/gum/skate-pass.gif" width="600" alt="Selecting a skate value with gum" />

#### Choose packages to uninstall

List all packages installed by your package manager (we'll use `brew`) and
choose which packages to uninstall.

```bash
brew list | gum choose --no-limit | xargs brew uninstall
```

#### Choose branches to delete

List all branches and choose which branches to delete.

```bash
git branch | cut -c 3- | gum choose --no-limit | xargs git branch -D
```

#### Choose pull request to checkout

List all PRs for the current GitHub repository and checkout the chosen PR (using [`gh`](https://cli.github.com/)).

```bash
gh pr list | cut -f1,2 | gum choose | cut -f1 | xargs gh pr checkout
```

#### Pick command from shell history

Pick a previously executed command from your shell history to execute, copy,
edit, etc...

```bash
gum filter < $HISTFILE --height 20
```

#### Sudo password input

See visual feedback when entering password with masked characters with `gum
input --password`.

```bash
alias please="gum input --password | sudo -nS"
```

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.social/@charmcli)
* [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
