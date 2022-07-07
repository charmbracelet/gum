# Gum

Tasty Bubble Gum for your shell.

Gum is a command-line tool (`gum`) that gives you the power of
[Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) without needing to write
any Go code.

The following example is generated entirely from a [single bash script](./examples/demo.sh).

<img src="https://stuff.charm.sh/gum/gum.gif" width="900" alt="Shell running the Gum examples/demo.sh script">

## Interaction

#### Input
Prompt your users for input with a simple command.

```bash
gum input > answer.text
```

#### Search

Allow your users to filter through a list of options by fuzzy searching.

```bash
echo Strawberry >> flavors.text
echo Banana >> flavors.text
echo Cherry >> flavors.text
cat flavors.text | gum search > selection.text
```

#### Loading

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

```
⣽ Buying Bubble Gum...
```


## Styling

Pretty print any string with any layout with one command.

```bash
gum style \
	--foreground "#FF06B7" --border "double" --align "center" \
	--width 50 --margin 2 --padding "2 4" \
	"Bubble Gum (1¢)" "So sweet and so fresh\!"
```
                                                        
Result:

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

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

* [Twitter](https://twitter.com/charmcli)
* [The Fediverse](https://mastodon.technology/@charm)
* [Slack](https://charm.sh/slack)

## License

[MIT](https://github.com/charmbracelet/seashell/raw/main/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
