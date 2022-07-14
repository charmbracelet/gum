# Gum Format

Gum format allows you to format different text into human readable output.

Four different parse-able formats exist:

1. [Markdown](#markdown)
2. [Code](#code)
3. [Template](#template)
3. [Emoji](#emoji)

## Markdown

Render any input as markdown text. This uses
[Glamour](https://github.com/charmbracelet/glamour) behind the scenes.

You can pass input as lines directly as arguments to the command invocation or
pass markdown over `stdin`.

```bash
gum format --type markdown < README.md
# Or, directly as arguments (useful for quick lists)
gum format --type markdown -- "# Gum Formats" "- Markdown" "- Code" "- Template" "- Emoji"
```

## Code

Render any code snippet with syntax highlighting.
[Glamour](https://github.com/charmbracelet/glamour), which uses
[Chroma](https://github.com/alecthomas/chroma) under the hood, handles styling.

Similarly to the `markdown` format, `code` can take input over `stdin`.

```bash
cat options.go | gum format --type code
```

## Template

Render styled input from a string template. Templating is handled by
[Termenv](https://github.com/muesli/termenv).

```bash
gum format --type template '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}'
# Or, via stdin
echo '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' | gum format --type template 
```

## Emoji

Parse and render emojis from their matching `:name:`s. Powered by
[Glamour](https://github.com/charmbracelet/glamour) and [Goldmark
Emoji](https://github.com/yuin/goldmark-emoji)

```bash
gum format --type emoji 'I :heart: Bubble Gum :candy:'
# You know the drill, also via stdin
echo 'I :heart: Bubble Gum :candy:' | gum format --type emoji 
```
