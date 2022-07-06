# Soda Pop

Command line utilities to make your command-line `pop`.

Soda Pop provides utilities to help you create shell scripts make your user's
shells `pop`. Powered by [Bubbles](https://github.com/charmbracelet/bubbles)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Interaction

#### Input
Prompt your users for input with a simple command.

```bash
pop input > answer.text
```

#### Search

Allow your users to filter through a list of options by fuzzy searching.

```bash
echo Gastropoda >> options.text
echo Bivalvia >> options.text
echo Polyplacophora >> options.text
cat options.text | pop search > selection.text
```

#### Loading

Display a progress bar while loading.

```bash
pop loading --time 5s
```

#### Spinners

Display a spinner while taking some running action. We specify the command to
run while showing the spinner, the spinner will automatically stop after the
command exits.

```bash
pop spin --type dot --title "Purchasing Seashells" -- sleep 5
```

```
⣽ Purchasing Seashells
```


## Styling

Pretty print any string with any layout with one command.

```bash
pop style \
		--foreground "#FF06B7" --border "double" \
		--margin 2 --padding "2 4" --width 50 \
		"And oh gosh, how delicious the fabulous frizzy frobscottle was!
```

```bash
pop style "She sells sea shells by the sea shore." --width 50 --height 5
  --align center --padding 2 --margin "1 2" --border double
```
                                                        
Result:

```
╔══════════════════════════════════════════════════╗
║                                                  ║
║                                                  ║
║    And oh gosh, how delicious the fabulous       ║
║    frizzy frobscottle was!                       ║
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
