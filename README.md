# Gum

Gum is a collection of command-line utilities that make your shell scripts a
little more glamorous. It gives you the power of
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) without needing to write
any Go code.

The following example is running from a [single Bash script](./examples/demo.sh).

<img src="https://stuff.charm.sh/gum/gum.gif" width="900" alt="Shell running the Gum examples/demo.sh script">

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
gum join "$(gum style --padding "1 5" --border double "I love Bubble Gum")" "$(gum style --padding "1 5" --border double "<3")"
```

```
╔═══════════════════════════╗╔════════════╗
║                           ║║            ║
║     I love Bubble Gum     ║║     <3     ║
║                           ║║            ║
╚═══════════════════════════╝╚════════════╝
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
