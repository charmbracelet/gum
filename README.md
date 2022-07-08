# Gum

<p>
    <a href="https://github.com/charmbracelet/gum/releases"><img src="https://img.shields.io/github/release/charmbracelet/gum.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/gum?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/gum/actions"><img src="https://github.com/charmbracelet/gum/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://nightly.link/charmbracelet/gum/workflows/nightly/main"><img src="https://shields.io/badge/-Nightly%20Builds-orange?logo=hackthebox&logoColor=fff&style=appveyor"/></a>
</p>

Gum is a collection of command-line utilities that make your shell scripts a
little more glamorous. It gives you the power of
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) without needing to write
any Go code.

The following example is running from a [single Bash script](./examples/demo.sh).

<img src="https://stuff.charm.sh/gum/gum.gif" width="900" alt="Shell running the Gum examples/demo.sh script">

## Installation

Use a package manager:

```bash
# macOS or Linux
brew tap charmbracelet/tap && brew install charmbracelet/tap/gum

# Arch Linux
pacman -S gum

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

## Interaction

#### Input
Prompt your users for input with a simple command.

```bash
gum input > answer.text
```

#### Write

Prompt your users to write some multi-line text.

```bash
gum write > story.text
```

#### Search

Allow your users to filter through a list of options by fuzzy searching.

```bash
echo Strawberry >> flavors.text
echo Banana >> flavors.text
echo Cherry >> flavors.text
cat flavors.text | gum search > selection.text
```

#### Progress

Display a progress bar while loading. The following command will display a
progress bar and increment the progress by 10% every 1 second. Thus, taking 10
seconds to complete the progress bar.

```bash
gum progress --increment 0.1 --interval 1s
```

#### Spinners

Display a spinner while taking some running action. We specify the command to
run while showing the spinner, the spinner will automatically stop after the
command exits.

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

## Styling and Layout

#### Style
Pretty print any string with any layout with one command.

```bash
gum style \
	--foreground "#FF06B7" --border "double" --align "center" \
	--width 50 --margin "1 2" --padding "2 4" \
	"Bubble Gum (1¢)" "So sweet and so fresh\!"
```

```
                                                        
  ╔══════════════════════════════════════════════════╗  
  ║                                                  ║  
  ║                                                  ║  
  ║                 Bubble Gum (1¢)                  ║  
  ║              So sweet and so fresh!              ║  
  ║                                                  ║  
  ║                                                  ║  
  ╚══════════════════════════════════════════════════╝  

```

#### Join

Combine text vertically or horizontally with a single command, use this command
with `gum style` to build layouts and pretty output.

Note: It's important to wrap the output of `gum style` in quotes to ensure new
lines (`\n`) are part of a single argument passed to the `join` command.

```bash
I=$(gum style --padding "1 5" --border double "I")
LOVE=$(gum style --padding "1 4" --border double "LOVE")
BUBBLE=$(gum style --padding "1 8" --border double "Bubble")
GUM=$(gum style --padding "1 5" --border double "Gum")

I_LOVE=$(gum join "$I" "$LOVE")
BUBBLE_GUM=$(gum join "$BUBBLE" "$GUM")
gum join --align center --vertical "$I_LOVE" "$BUBBLE_GUM"
```

```
      ╔═══════════╗╔════════════╗      
      ║           ║║            ║      
      ║     I     ║║    LOVE    ║      
      ║           ║║            ║      
      ╚═══════════╝╚════════════╝      
╔══════════════════════╗╔═════════════╗
║                      ║║             ║
║        Bubble        ║║     Gum     ║
║                      ║║             ║
╚══════════════════════╝╚═════════════╝
```


## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.technology/@charm)
* [Slack](https://charm.sh/slack)

## License

[MIT](https://github.com/charmbracelet/seashell/raw/main/LICENSE)

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
