# Gum

<p>
    <a href="https://github.com/charmbracelet/gum/releases"><img src="https://img.shields.io/github/release/charmbracelet/gum.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/gum?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/gum/actions"><img src="https://github.com/charmbracelet/gum/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Gum is a collection of command-line utilities that make your shell scripts a
little more glamorous. It gives you the power of
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) without needing to write
any Go code.

```bash
# Prompt users for input
NAME=$(gum input --placeholder "What is your name?")

# Style some text
gum style --foreground 212 --padding "1 4" \
	--border double --border-foreground 57 \
	"Nice to meet you, $NAME."

# Do some work while spinning
gum spin --title "Taking a nap..." --color 212 -- sleep 5

# Fuzzy find a file or directory
find . -type f | gum filter
```

The following example is running from a [single bash script](./examples/demo.sh).

<img src="https://stuff.charm.sh/gum/gum.gif" width="900" alt="Shell running the Gum examples/demo.sh script">

## Installation

Use a package manager:

```bash
# macOS or Linux
brew tap charmbracelet/tap && brew install charmbracelet/tap/gum

# Arch Linux (btw)
pacman -S gum

# Nix
nix-env -iA nixpkgs.gum

# Debian/Ubuntu
echo 'deb [trusted=yes] https://repo.charm.sh/apt/ /' \
    | sudo tee /etc/apt/sources.list.d/charm.list
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

## Interaction

#### Input
Prompt your users for input with a simple command.

```bash
gum input > answer.text
```

<img src="https://stuff.charm.sh/gum/input.gif" width="350" alt="Shell running gum input typing I love bubble gum <3">

#### Write

Prompt your users to write some multi-line text.

```bash
gum write > story.text
```

<img src="https://stuff.charm.sh/gum/write.gif" width="600" alt="Shell running gum write typing My favorite flavors are:">

#### Filter

Allow your users to filter through a list of options by fuzzy searching.

```bash
echo Strawberry >> flavors.text
echo Banana >> flavors.text
echo Cherry >> flavors.text
cat flavors.text | gum filter > selection.text
```

<img src="https://stuff.charm.sh/gum/filter.gif" width="600" alt="Shell running gum filter on different bubble gum flavors">

#### Choose

Ask your users to choose an option from a list of choices.

```bash
echo "Pick a card, any card..."
CARD=$(gum choose --height 15 {{A,K,Q,J},{10..2}}" "{♠,♥,♣,♦})
echo "Was your card the $CARD?"
```

<img src="https://stuff.charm.sh/gum/choose.gif" width="600" alt="Shell running gum choose on a deck of cards, picking the Ace of Hearts">

You can also set a limit on the number of items to choose with the `--limit` flag.

```bash
echo "Pick your top 5 songs."
cat songs.txt | gum choose --limit 5
```

Or, allow any number of selections with the `--no-limit` flag.

```bash
echo "What do you need from the grocery store?"
cat foods.txt | gum choose --no-limit
```

#### Progress

Display a progress bar while loading. The following command will display a
progress bar and increment the progress by 10% every 1 second. Thus, taking 10
seconds to complete the progress bar.

```bash
gum progress --increment 0.1 --interval 1s
```

<img src="https://stuff.charm.sh/gum/progress.gif" width="600" alt="Shell running gum progress">

#### Spinners

Display a spinner while taking some running action. We specify the command to
run while showing the spinner, the spinner will automatically stop after the
command exits.

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

<img src="https://stuff.charm.sh/gum/spin.gif?cache=1" width="600" alt="Shell running gum spin while sleeping for 5 seconds">

## Styling and Layout

#### Style
Pretty print any string with any layout with one command.

```bash
gum style \
	--foreground 212 --border-foreground 212 --border double \
	--align center --width 50 --margin "1 2" --padding "2 4" \
	'Bubble Gum (1¢)' 'So sweet and so fresh!'
```

<img src="https://stuff.charm.sh/gum/style.png?cache=1" width="600" alt="Bubble Gum, So sweet and so fresh!">

#### Join

Combine text vertically or horizontally with a single command, use this command
with `gum style` to build layouts and pretty output.

Note: It's important to wrap the output of `gum style` in quotes to ensure new
lines (`\n`) are part of a single argument passed to the `join` command.

```bash
I=$(gum style --padding "1 5" --border double --border-foreground 212 "I")
LOVE=$(gum style --padding "1 4" --border double --border-foreground 57 "LOVE")
BUBBLE=$(gum style --padding "1 8" --border double --border-foreground 255 "Bubble")
GUM=$(gum style --padding "1 5" --border double --border-foreground 240 "Gum")

I_LOVE=$(gum join "$I" "$LOVE")
BUBBLE_GUM=$(gum join "$BUBBLE" "$GUM")
gum join --align center --vertical "$I_LOVE" "$BUBBLE_GUM"
```

<img src="https://stuff.charm.sh/gum/join.png" width="600" alt="I LOVE Bubble Gum written out in four boxes with double borders around them.">

#### Format

Format allows you to take some text and stylize it. `gum format` can parse
markdown, code, template strings, and emoji strings.

For example, we can output a markdown list with the following command:

```bash
gum format -- "# Gum Formats" "- Markdown" "- Code" "- Template" "- Emoji"
echo "# Gum Formats\n- Markdown\n- Code\n- Template\n- Emoji" | gum format
```

<img src="https://stuff.charm.sh/gum/format-markdown.png" width="250" alt="Format markdown displaying a markdown list with content of different possible formats">

Apply syntax highlighting to code with `gum format -t code`:

```bash
cat main.go | gum format -t code
```

<img src="https://stuff.charm.sh/gum/format-code.png" width="400" alt="Format code command displaying a simple Hello, world! Go program with syntax highlighting.">

Render a template string with `gum format -t template`:

```bash
echo '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' \
    | gum format -t template
```

<img src="https://stuff.charm.sh/gum/format-template.png" width="400" alt="Format template command displaying Tasty Bubble Gum in different styles">

Render some emojis with `:emoji:` syntax with `gum format -t emoji`.

```bash
echo 'I :heart: Bubble Gum :candy:' | gum format -t emoji
```

<img src="https://stuff.charm.sh/gum/format-emoji.png" width="400" alt="Format emoji command displaying I :heart: Bubble Gum :candy:.">


## Examples

See the [examples](./examples/) directory for more real world use cases.

How to use `gum` in your daily workflows:

#### Write a commit message

Prompt for user input to write git commit messages with a short summary and
longer details with `gum input` and `gum write`.

Bonus points if you use `gum filter` with the [Conventional Commits
Specification](https://www.conventionalcommits.org/en/v1.0.0/#summary) as a
prefix for your commit message.

```bash
git commit -m "$(gum input --width 50 --placeholder "Summary of changes")" \
           -m "$(gum write --width 80 --placeholder "Details of changes")"
```

<img src="https://stuff.charm.sh/gum/commit.gif" width="600" alt="Running the ./examples/commit.sh script to commit to git">

#### Open files in your `$EDITOR`

By default `gum filter` will display a list of all files (searched recursively)
through your current directory, it has some sensible ignored defaults (`.git`,
`node_modules`). You can use this to pick a file and open it in your `$EDITOR`.

```bash
$EDITOR $(gum filter)
```

#### Connect to a TMUX session

Pick from a running `TMUX` session and attach to it if not inside `TMUX` or
switch your client to the session if already attached to a session.

```bash
SESSION=$(tmux list-sessions -F \#S | gum filter --placeholder "Pick session...")
tmux switch-client -t $SESSION || tmux attach -t $SESSION
```

<img src="https://stuff.charm.sh/gum/tmux-session-picker.png" width="250" alt="Picking a tmux session with gum filter">

#### Pick commit hash from history

Filter through your git history searching for commit messages and copy the
commit hash of the selected commit.

```bash
git log --oneline | gum filter | cut -d' ' -f1 # | copy
```

<img src="https://stuff.charm.sh/gum/commit-picker.png" width="500" alt="Picking a commit with gum filter">

#### Choose packages to uninstall

List all packages installed by your package manager (we'll use `brew`) and
choose which packages to uninstall.

```bash
brew list | gum choose --no-limit | xargs brew uninstall
```

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.technology/@charm)
* [Slack](https://charm.sh/slack)

## License

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
