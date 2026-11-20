# Gum Format

> **语言**：[English](https://github.com/charmbracelet/gum/blob/main/format/README.md) | 简体中文（本文档）
>
> 本文档为 gum `format` 子命令文档的简体中文译本。所有命令、选项、模板语法与脚本代码均保持原文。

Gum format 可以把各种文本格式化为人类可读的输出。

共有四种可解析的格式：

1. [Markdown](#markdown)
2. [Code](#code)
3. [Template](#template)
4. [Emoji](#emoji)

## Markdown

把任意输入渲染为 Markdown 文本。底层使用
[Glamour](https://github.com/charmbracelet/glamour)。

既可以把输入作为参数逐行直接传给命令，
也可以通过 `stdin` 传入 Markdown。

```bash
gum format --type markdown < README.md
# Or, directly as arguments (useful for quick lists)
gum format --type markdown -- "# Gum Formats" "- Markdown" "- Code" "- Template" "- Emoji"
```

## Code

渲染任意代码片段，并带上语法高亮。
样式由 [Glamour](https://github.com/charmbracelet/glamour) 处理，它底层使用
[Chroma](https://github.com/alecthomas/chroma)。

与 `markdown` 格式类似，`code` 也可以从 `stdin` 读取输入。

```bash
cat options.go | gum format --type code
```

## Template

根据字符串模板渲染带样式的输入。模板由
[Termenv](https://github.com/muesli/termenv) 处理。

```bash
gum format --type template '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}'
# Or, via stdin
echo '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' | gum format --type template
```

## Emoji

解析 `:name:` 形式的写法，并渲染出对应的 emoji。由
[Glamour](https://github.com/charmbracelet/glamour) 与 [Goldmark
Emoji](https://github.com/yuin/goldmark-emoji) 提供支持

```bash
gum format --type emoji 'I :heart: Bubble Gum :candy:'
# You know the drill, also via stdin
echo 'I :heart: Bubble Gum :candy:' | gum format --type emoji
```

## Tables

表格使用 [Glamour](https://github.com/charmbracelet/glamour) 渲染。

| 泡泡糖口味 | 价格  |
| ----------------- | ----- |
| Strawberry        | $0.99 |
| Cherry            | $0.50 |
| Banana            | $0.75 |
| Orange            | $0.25 |
| Lemon             | $0.50 |
| Lime              | $0.50 |
| Grape             | $0.50 |
| Watermelon        | $0.50 |
| Pineapple         | $0.50 |
| Blueberry         | $0.50 |
| Raspberry         | $0.50 |
| Cranberry         | $0.50 |
| Peach             | $0.50 |
| Apple             | $0.50 |
| Mango             | $0.50 |
| Pomegranate       | $0.50 |
| Coconut           | $0.50 |
| Cinnamon          | $0.50 |
