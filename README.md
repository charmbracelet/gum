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

<img src="https://stuff.charm.sh/gum/demo.gif" width="600" alt="Shell running the ./demo.sh script" />

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
gum write --placeholder "Details of this change"
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
DESCRIPTION=$(gum write --placeholder "Details of this change")

# Commit these changes
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
```

<img src="https://stuff.charm.sh/gum/commit_2.gif" width="600" alt="Running the ./examples/commit.sh script to commit to git" />

## Installation

Use a package manager:

```bash
# macOS or Linux
brew install gum

# Arch Linux (btw)
yay -S gum-bin

# Nix
nix-env -iA nixpkgs.gum

# Debian/Ubuntu
echo 'deb [trusted=yes] https://repo.charm.sh/apt/ /' | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install gum

# Fedora
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/charm.repo
sudo yum install gum
```

Or download it:

* [Packages][releases] are available in Debian and RPM formats
* [Binaries][releases] are available for Linux, macOS, and Windows

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

<img src="https://stuff.charm.sh/gum/customization.gif" width="600" alt="Gum input displaying most customization options" />

## Interaction

#### Input

Prompt for input with a simple command.

```bash
gum input > answer.text
```

<img src="https://stuff.charm.sh/gum/input_1.gif" width="600" alt="Shell running gum input typing Not much, you?" />

#### Write

Prompt for some multi-line text.

```bash
gum write > story.text
```

<img src="https://stuff.charm.sh/gum/write.gif" width="600" alt="Shell running gum write typing a story" />

#### Filter

Use fuzzy matching to filter a list of values:

```bash
echo Strawberry >> flavors.text
echo Banana >> flavors.text
echo Cherry >> flavors.text
cat flavors.text | gum filter > selection.text
```

<img src="https://stuff.charm.sh/gum/filter.gif" width="600" alt="Shell running gum filter on different bubble gum flavors" />

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


#### Spin

Display a spinner while running a script or command. The spinner will
automatically stop after the given command exits.

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

<img src="https://stuff.charm.sh/gum/spin.gif" width="600" alt="Shell running gum spin while sleeping for 5 seconds" />

Available spinner types include: `line`, `dot`, `minidot`, `jump`, `pulse`, `points`, `globe`, `moon`, `monkey`, `meter`, `hamburger`.

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
           -m "$(gum write --width 80 --placeholder "Details of changes")"
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

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.technology/@charm)
* [Slack](https://charm.sh/slack)

## License

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
