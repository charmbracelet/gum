# Sea Shell

Command line utilities to make delightful shell scripts.

Sea Shell provides utilities to help you create shell scripts that will delight
your users. Powered by [Bubbles](https://github.com/charmbracelet/bubbles) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Interaction

#### Input
Prompt your users for input with a simple command.

```bash
sea input > answer.text
```

#### Search

Allow your users to filter through a list of options by fuzzy searching.

```bash
echo Gastropoda >> options.text
echo Bivalvia >> options.text
echo Polyplacophora >> options.text
cat options.text | grg search > selection.text
```

#### Loading

Display a progress bar while loading.

```bash
sea loading --time 5s
```

#### Spinners

Display a spinner while taking some running action. We specify the command to
run while showing the spinner, the spinner will automatically stop after the
command exits.

```bash
sea spin --type dot --title "Purchasing Seashells" -- sleep 5
```

```
⣽ Purchasing Seashells
```


## Styling

Pretty print any string with any layout with one command.

```bash
sea style "She sells sea shells by the sea shore." --width 50 --height 5
  --align center --padding 2 --margin "1 2" --border double
```
                                                        
Result:

```
                                                        
  ╔══════════════════════════════════════════════════╗  
  ║                                                  ║  
  ║                                                  ║  
  ║      She sells sea shells by the sea shore.      ║  
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
