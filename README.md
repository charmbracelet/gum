# Gum

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

<img alt="Shell running the ./demo.sh script" width="600" src="https://vhs.charm.sh/vhs-1qY57RrQlXCuydsEgDp68G.gif">

The above example is running from a single shell script ([source](./examples/demo.sh)).

## Tutorial

Gum provides highly configurable, ready-to-use utilities to help you write
useful shell scripts and dotfiles aliases with just a few lines of code.
Let's build a simple script to help you write
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)
for your dotfiles.

Ask for the commit type with gum choose:

```bash
gum choose "fix" "feat" "docs" "style" "refactor" "test" "chore" "revert"
```

> [!NOTE]
> This command itself will print to stdout which is not all that useful. To make use of the command later on you can save the stdout to a `$VARIABLE` or `file.txt`.

Prompt for the scope of these changes:

```bash
gum input --placeholder "scope"
```

Prompt for the summary and description of changes:

```bash
gum input --value "$TYPE$SCOPE: " --placeholder "Summary of this change"
gum write --placeholder "Details of this change"
```

Confirm before committing:

```bash
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
```

Check out the [complete example](https://github.com/charmbracelet/gum/blob/main/examples/commit.sh) for combining these commands in a single script.

<img alt="Running the ./examples/commit.sh script to commit to git" width="600" src="https://vhs.charm.sh/vhs-7rRq3LsEuJVwhwr0xf6Er7.gif">

## Installation

Use a package manager:

```bash
# macOS or Linux
brew install gum

# Arch Linux (btw)
pacman -S gum

# Nix
nix-env -iA nixpkgs.gum

# Flox
flox install gum

# Windows (via WinGet or Scoop)
winget install charmbracelet.gum
scoop install charm-gum
```

<details>
<summary>Debian/Ubuntu</summary>

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install gum
```

</details>

<details>
<summary>Fedora/RHEL/OpenSuse</summary>

```bash
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
sudo rpm --import https://repo.charm.sh/yum/gpg.key

# yum
sudo yum install gum

# zypper
sudo zypper refresh
sudo zypper install gum
```

</details>

<details>
<summary>FreeBSD</summary>

```bash
# packages
sudo pkg install gum

# ports
cd /usr/ports/devel/gum && sudo make install clean
```

</details>

Or download it:

- [Packages][releases] are available in Debian, RPM, and Alpine formats
- [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

Or just install it with `go`:

```bash
go install github.com/charmbracelet/gum@latest
```

[releases]: https://github.com/charmbracelet/gum/releases

## Commands

- [`choose`](#choose): Choose an option from a list of choices
- [`confirm`](#confirm): Ask a user to confirm an action
- [`file`](#file): Pick a file from a folder
- [`filter`](#filter): Filter items from a list
- [`format`](#format): Format a string using a template
- [`input`](#input): Prompt for some input
- [`join`](#join): Join text vertically or horizontally
- [`pager`](#pager): Scroll through a file
- [`spin`](#spin): Display spinner while running a command
- [`style`](#style): Apply coloring, borders, spacing to text
- [`table`](#table): Render a table of data
- [`write`](#write): Prompt for long-form text
- [`log`](#log): Log messages to output

## Customization

You can customize `gum` options and styles with `--flags` and `$ENVIRONMENT_VARIABLES`.
See `gum <command> --help` for a full view of each command's customization and configuration options.

Customize with `--flags`:

```bash

gum input --cursor.foreground "#FF0" \
          --prompt.foreground "#0FF" \
          --placeholder "What's up?" \
          --prompt "* " \
          --width 80 \
          --value "Not much, hby?"
```

Customize with `ENVIRONMENT_VARIABLES`:

```bash
export GUM_INPUT_CURSOR_FOREGROUND="#FF0"
export GUM_INPUT_PROMPT_FOREGROUND="#0FF"
export GUM_INPUT_PLACEHOLDER="What's up?"
export GUM_INPUT_PROMPT="* "
export GUM_INPUT_WIDTH=80

# --flags can override values set with environment
gum input
```

<img alt="Gum input displaying most customization options" width="600" src="https://vhs.charm.sh/vhs-5zb9DlQYA70aL9ZpYLTwKv.gif">

## Input

Prompt for input with a simple command.

```bash
gum input > answer.txt
gum input --password > password.txt
```

<img src="https://vhs.charm.sh/vhs-1nScrStFI3BMlCp5yrLtyg.gif" width="600" alt="Shell running gum input typing Not much, you?" />

## Write

Prompt for some multi-line text (`ctrl+d` to complete text entry).

```bash
gum write > story.txt
```

<img src="https://vhs.charm.sh/vhs-7abdKKrUEukgx9aJj8O5GX.gif" width="600" alt="Shell running gum write typing a story" />

## Filter

Filter a list of values with fuzzy matching:

```bash
echo Strawberry >> flavors.txt
echo Banana >> flavors.txt
echo Cherry >> flavors.txt
gum filter < flavors.txt > selection.txt
```

<img src="https://vhs.charm.sh/vhs-61euOQtKPtQVD7nDpHQhzr.gif" width="600" alt="Shell running gum filter on different bubble gum flavors" />

Select multiple options with the `--limit` flag or `--no-limit` flag. Use `tab` or `ctrl+space` to select, `enter` to confirm.

```bash
cat flavors.txt | gum filter --limit 2
cat flavors.txt | gum filter --no-limit
```

## Choose

Choose an option from a list of choices.

```bash
echo "Pick a card, any card..."
CARD=$(gum choose --height 15 {{A,K,Q,J},{10..2}}" "{♠,♥,♣,♦})
echo "Was your card the $CARD?"
```

You can also select multiple items with the `--limit` or `--no-limit` flag, which determines
the maximum of items that can be chosen.

```bash
cat songs.txt | gum choose --limit 5
cat foods.txt | gum choose --no-limit --header "Grocery Shopping"
```

<img src="https://vhs.charm.sh/vhs-3zV1LvofA6Cbn5vBu1NHHl.gif" width="600" alt="Shell running gum choose with numbers and gum flavors" />

## Confirm

Confirm whether to perform an action. Exits with code `0` (affirmative) or `1`
(negative) depending on selection.

```bash
gum confirm && rm file.txt || echo "File not removed"
```

<img src="https://vhs.charm.sh/vhs-3xRFvbeQ4lqGerbHY7y3q2.gif" width="600" alt="Shell running gum confirm" />

## File

Prompt the user to select a file from the file tree.

```bash
$EDITOR $(gum file $HOME)
```

<img src="https://vhs.charm.sh/vhs-2RMRqmnOPneneIgVJJ3mI1.gif" width="600" alt="Shell running gum file" />

## Pager

Scroll through a long document with line numbers and a fully customizable viewport.

```bash
gum pager < README.md
```

<img src="https://vhs.charm.sh/vhs-3iMDpgOLmbYr0jrYEGbk7p.gif" width="600" alt="Shell running gum pager" />

## Spin

Display a spinner while running a script or command. The spinner will
automatically stop after the given command exits.

To view or pipe the command's output, use the `--show-output` flag.

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

<img src="https://vhs.charm.sh/vhs-3YFswCmoY4o3Q7MyzWl6sS.gif" width="600" alt="Shell running gum spin while sleeping for 5 seconds" />

Available spinner types include: `line`, `dot`, `minidot`, `jump`, `pulse`, `points`, `globe`, `moon`, `monkey`, `meter`, `hamburger`.

## Table

Select a row from some tabular data.

```bash
gum table < flavors.csv | cut -d ',' -f 1
```

<!-- <img src="https://stuff.charm.sh/gum/table.gif" width="600" alt="Shell running gum table" /> -->

## Style

Pretty print any string with any layout with one command.

```bash
gum style \
	--foreground 212 --border-foreground 212 --border double \
	--align center --width 50 --margin "1 2" --padding "2 4" \
	'Bubble Gum (1¢)' 'So sweet and so fresh!'
```

<img src="https://github.com/charmbracelet/gum/assets/42545625/67468acf-b3e0-4e78-bd89-360739eb44fa" width="600" alt="Bubble Gum, So sweet and so fresh!" />

## Join

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

<img src="https://github.com/charmbracelet/gum/assets/42545625/68f7a25d-b495-48dd-982a-cee0c8ea5786" width="600" alt="I LOVE Bubble Gum written out in four boxes with double borders around them." />

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

<img src="https://github.com/charmbracelet/gum/assets/42545625/5cfbb0c8-0022-460d-841b-fec37527ca66" width="300" alt="Running gum format for different types of formats" />

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

# Include a timestamp.
gum log --time rfc822 --level error "Unable to create file."
```

See the Go [`time` package](https://pkg.go.dev/time#pkg-constants) for acceptable `--time` formats.

See [`charmbracelet/log`](https://github.com/charmbracelet/log) for more usage.

<img src="https://vhs.charm.sh/vhs-6jupuFM0s2fXiUrBE0I1vU.gif" width="600" alt="Running gum log with debug and error levels" />

## Examples

How to use `gum` in your daily workflows:

See the [examples](./examples/) directory for more real world use cases.

- Write a commit message:

```bash
git commit -m "$(gum input --width 50 --placeholder "Summary of changes")" \
           -m "$(gum write --width 80 --placeholder "Details of changes")"
```

- Open files in your `$EDITOR`

```bash
$EDITOR $(gum filter)
```

- Connect to a `tmux` session

```bash
SESSION=$(tmux list-sessions -F \#S | gum filter --placeholder "Pick session...")
tmux switch-client -t "$SESSION" || tmux attach -t "$SESSION"
```

- Pick a commit hash from `git` history

```bash
git log --oneline | gum filter | cut -d' ' -f1 # | copy
```

- Simple [`skate`](https://github.com/charmbracelet/skate) password selector.

```
skate list -k | gum filter | xargs skate get
```

- Uninstall packages

```bash
brew list | gum choose --no-limit | xargs brew uninstall
```

- Clean up `git` branches

```bash
git branch | cut -c 3- | gum choose --no-limit | xargs git branch -D
```

- Checkout GitHub pull requests with [`gh`](https://cli.github.com/)

```bash
gh pr list | cut -f1,2 | gum choose | cut -f1 | xargs gh pr checkout
```

- Copy command from shell history

```bash
gum filter < $HISTFILE --height 20
```

- `sudo` replacement

```bash
alias please="gum input --password | sudo -nS"
```

## Contributing

See [contributing][contribute].

[contribute]: https://github.com/charmbracelet/gum/contribute

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
