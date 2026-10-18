# Gum

> **语言**：[English](https://github.com/charmbracelet/gum/blob/main/README.md) | 简体中文（本文档）
>
> 本文档为 gum 官方 README 的简体中文译本。所有命令、选项、环境变量名与脚本代码均保持原文。

<p>
    <a href="https://stuff.charm.sh/gum/nutritional-information.png" target="_blank"><img src="https://stuff.charm.sh/gum/gum.png" alt="Gum Image" width="450" /></a>
    <br><br>
    <a href="https://github.com/charmbracelet/gum/releases"><img src="https://img.shields.io/github/release/charmbracelet/gum.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/gum?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
    <a href="https://github.com/charmbracelet/gum/actions"><img src="https://github.com/charmbracelet/gum/workflows/build/badge.svg" alt="Build Status"></a>
</p>

一个让 shell 脚本变得赏心悦目的工具。你无需编写任何 Go 代码，
就能在脚本与别名中借助
[Bubbles](https://github.com/charmbracelet/bubbles) 与 [Lip
Gloss](https://github.com/charmbracelet/lipgloss) 的能力！

<img alt="Shell running the ./demo.sh script" width="600" src="https://vhs.charm.sh/vhs-1qY57RrQlXCuydsEgDp68G.gif">

上面的示例完全由一个 shell 脚本驱动（[源码](./examples/demo.sh)）。

## 教程

Gum 提供了高度可配置、开箱即用的实用工具，让你只需
几行代码就能写出有用的 shell 脚本与 dotfile 别名。
我们来写一个简单的脚本，帮你为 dotfiles 撰写
[Conventional Commits（约定式提交）](https://www.conventionalcommits.org/zh-hans/v1.0.0/#%E6%A6%82%E8%A6%81)。

用 gum choose 询问提交类型：

```bash
gum choose "fix" "feat" "docs" "style" "refactor" "test" "chore" "revert"
```

> [!NOTE]
> 这条命令本身只是把结果打印到 stdout，用处不大。要想在后面使用这条命令的结果，你可以把 stdout 保存到 `$VARIABLE` 或 `file.txt`。

询问这些改动的影响范围（scope）：

```bash
gum input --placeholder "scope"
```

询问改动的摘要与描述：

```bash
gum input --value "$TYPE$SCOPE: " --placeholder "Summary of this change"
gum write --placeholder "Details of this change"
```

在提交前进行确认：

```bash
gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
```

把上述命令组合成单个脚本的[完整示例](https://github.com/charmbracelet/gum/blob/main/examples/commit.sh)在此。

<img alt="Running the ./examples/commit.sh script to commit to git" width="600" src="https://vhs.charm.sh/vhs-7rRq3LsEuJVwhwr0xf6Er7.gif">

## 安装

使用包管理器：

```bash
# macOS 或 Linux
brew install gum

# Arch Linux (btw)
pacman -S gum

# Fedora 或 EPEL 10
dnf install gum

# Nix
nix-env -iA nixpkgs.gum

# Flox
flox install gum

# Windows（通过 WinGet 或 Scoop）
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

或者直接下载：

- [软件包][releases]提供 Debian、RPM 与 Alpine 格式
- [二进制文件][releases]提供 Linux、macOS、Windows、FreeBSD、OpenBSD 与 NetBSD 版本

或者直接用 `go` 安装：

```bash
go install charm.land/gum/v2@latest
```

[releases]: https://github.com/charmbracelet/gum/releases

## 命令

- [`choose`](#choose选择)：从一组选项中选择一项
- [`confirm`](#confirm确认)：请用户确认是否执行某个操作
- [`file`](#file文件)：从文件夹中选择一个文件
- [`filter`](#filter筛选)：从列表中筛选条目
- [`format`](#format格式化)：用模板格式化字符串
- [`input`](#input输入)：提示输入一些内容
- [`join`](#join拼接)：纵向或横向拼接文本
- [`pager`](#pager分页浏览)：滚动浏览文件
- [`spin`](#spin转圈提示)：在运行命令时显示 spinner
- [`style`](#style样式)：为文本应用配色、边框与间距
- [`table`](#table表格)：渲染数据表格
- [`write`](#write长文本输入)：提示输入长文本
- [`log`](#log日志)：把日志消息输出到终端

## 自定义

你可以通过 `--flags` 与 `$ENVIRONMENT_VARIABLES` 来自定义 `gum` 的选项与样式。
运行 `gum <command> --help` 可查看每条命令完整的自定义与配置项。

用 `--flags` 自定义：

```bash

gum input --cursor.foreground "#FF0" \
          --prompt.foreground "#0FF" \
          --placeholder "What's up?" \
          --prompt "* " \
          --width 80 \
          --value "Not much, hby?"
```

用 `ENVIRONMENT_VARIABLES` 自定义：

```bash
export GUM_INPUT_CURSOR_FOREGROUND="#FF0"
export GUM_INPUT_PROMPT_FOREGROUND="#0FF"
export GUM_INPUT_PLACEHOLDER="What's up?"
export GUM_INPUT_PROMPT="* "
export GUM_INPUT_WIDTH=80

# --flags 可以覆盖通过环境变量设置的值
gum input
```

<img alt="Gum input displaying most customization options" width="600" src="https://vhs.charm.sh/vhs-5zb9DlQYA70aL9ZpYLTwKv.gif">

## Input（输入）

用一条简单的命令来提示输入。

```bash
gum input > answer.txt
gum input --password > password.txt
```

<img src="https://vhs.charm.sh/vhs-1nScrStFI3BMlCp5yrLtyg.gif" width="600" alt="Shell running gum input typing Not much, you?" />

## Write（长文本输入）

提示输入多行文本（`ctrl+d` 结束输入）。

```bash
gum write > story.txt
```

<img src="https://vhs.charm.sh/vhs-7abdKKrUEukgx9aJj8O5GX.gif" width="600" alt="Shell running gum write typing a story" />

## Filter（筛选）

用模糊匹配筛选一组值：

```bash
echo Strawberry >> flavors.txt
echo Banana >> flavors.txt
echo Cherry >> flavors.txt
gum filter < flavors.txt > selection.txt
```

<img src="https://vhs.charm.sh/vhs-61euOQtKPtQVD7nDpHQhzr.gif" width="600" alt="Shell running gum filter on different bubble gum flavors" />

用 `--limit` 或 `--no-limit` 标志选择多个选项。用 `tab` 或 `ctrl+space` 选择，用 `enter` 确认。

```bash
cat flavors.txt | gum filter --limit 2
cat flavors.txt | gum filter --no-limit
```

## Choose（选择）

从一组选项中选择一项。

```bash
echo "Pick a card, any card..."
CARD=$(gum choose --height 15 {{A,K,Q,J},{10..2}}" "{♠,♥,♣,♦})
echo "Was your card the $CARD?"
```

你也可以用 `--limit` 或 `--no-limit` 标志选择多个条目，它决定了
最多可以选择多少项。

```bash
cat songs.txt | gum choose --limit 5
cat foods.txt | gum choose --no-limit --header "Grocery Shopping"
```

<img src="https://vhs.charm.sh/vhs-3zV1LvofA6Cbn5vBu1NHHl.gif" width="600" alt="Shell running gum choose with numbers and gum flavors" />

## Confirm（确认）

确认是否执行某个操作。根据选择结果，以代码 `0`（肯定）或 `1`
（否定）退出。

```bash
gum confirm && rm file.txt || echo "File not removed"
```

<img src="https://vhs.charm.sh/vhs-3xRFvbeQ4lqGerbHY7y3q2.gif" width="600" alt="Shell running gum confirm" />

## File（文件）

提示用户从文件树中选择一个文件。

```bash
$EDITOR $(gum file $HOME)
```

<img src="https://vhs.charm.sh/vhs-2RMRqmnOPneneIgVJJ3mI1.gif" width="600" alt="Shell running gum file" />

## Pager（分页浏览）

带行号地滚动浏览一篇很长的文档，视口完全可定制。

```bash
gum pager < README.md
```

<img src="https://vhs.charm.sh/vhs-3iMDpgOLmbYr0jrYEGbk7p.gif" width="600" alt="Shell running gum pager" />

## Spin（转圈提示）

在运行脚本或命令时显示一个 spinner。该 spinner 会在
给定命令退出后自动停止。

如果想查看或通过管道获取该命令的输出，请使用 `--show-output` 标志。

```bash
gum spin --spinner dot --title "Buying Bubble Gum..." -- sleep 5
```

<img src="https://vhs.charm.sh/vhs-3YFswCmoY4o3Q7MyzWl6sS.gif" width="600" alt="Shell running gum spin while sleeping for 5 seconds" />

可用的 spinner 类型包括：`line`、`dot`、`minidot`、`jump`、`pulse`、`points`、`globe`、`moon`、`monkey`、`meter`、`hamburger`。

## Table（表格）

从一些表格化数据中选择一行。

```bash
gum table < flavors.csv | cut -d ',' -f 1
```

<!-- <img src="https://stuff.charm.sh/gum/table.gif" width="600" alt="Shell running gum table" /> -->

## Style（样式）

一条命令就能以任意版式漂亮地打印任意字符串。

```bash
gum style \
	--foreground 212 --border-foreground 212 --border double \
	--align center --width 50 --margin "1 2" --padding "2 4" \
	'Bubble Gum (1¢)' 'So sweet and so fresh!'
```

<img src="https://github.com/charmbracelet/gum/assets/42545625/67468acf-b3e0-4e78-bd89-360739eb44fa" width="600" alt="Bubble Gum, So sweet and so fresh!" />

## Join（拼接）

纵向或横向拼接文本。把这个命令与 `gum style` 配合使用，
即可搭建版式并输出漂亮的排版结果。

小提示：在 `join` 命令中把 `gum style` 的输出作为参数使用时，
请务必用引号把它包裹起来，以保留换行符（`\n`）。

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

## Format（格式化）

`format` 用于处理和格式化文本正文。`gum format` 可以解析 markdown、
模板字符串与具名 emoji。

```bash
# 格式化一些 markdown
gum format -- "# Gum Formats" "- Markdown" "- Code" "- Template" "- Emoji"
echo "# Gum Formats\n- Markdown\n- Code\n- Template\n- Emoji" | gum format

# 为一些代码做语法高亮
cat main.go | gum format -t code

# 用模板随心所欲地渲染文本
echo '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' \
    | gum format -t template

# 显示你喜欢的 emoji！
echo 'I :heart: Bubble Gum :candy:' | gum format -t emoji
```

有关模板辅助函数的更多信息，请参阅 [Termenv
文档](https://github.com/muesli/termenv#template-helpers)。完整的具名 emoji 列表
请参阅 [GitHub API](https://api.github.com/emojis)。

<img src="https://github.com/charmbracelet/gum/assets/42545625/5cfbb0c8-0022-460d-841b-fec37527ca66" width="300" alt="Running gum format for different types of formats" />

## Log（日志）

`log` 借助 [`charmbracelet/log`](https://github.com/charmbracelet/log) 库，
以不同的级别与样式把日志消息输出到终端。

```bash
# 输出一些调试信息。
gum log --structured --level debug "Creating file..." name file.txt
# DEBUG Creating file... name=file.txt

# 输出一些错误信息。
gum log --structured --level error "Unable to create file." name file.txt
# ERROR Unable to create file. name=temp.txt

# 带上时间戳。
gum log --time rfc822 --level error "Unable to create file."
```

`--time` 可接受的格式请参阅 Go 的 [`time` 包](https://pkg.go.dev/time#pkg-constants)。

更多用法请参阅 [`charmbracelet/log`](https://github.com/charmbracelet/log)。

<img src="https://vhs.charm.sh/vhs-6jupuFM0s2fXiUrBE0I1vU.gif" width="600" alt="Running gum log with debug and error levels" />

## 示例

如何在日常工作中使用 `gum`：

更多真实场景的用法请参见 [examples](./examples/) 目录。

- 写一条提交信息：

```bash
git commit -m "$(gum input --width 50 --placeholder "Summary of changes")" \
           -m "$(gum write --width 80 --placeholder "Details of changes")"
```

- 用你的 `$EDITOR` 打开文件

```bash
$EDITOR $(gum filter)
```

- 连接到某个 `tmux` 会话

```bash
SESSION=$(tmux list-sessions -F \#S | gum filter --placeholder "Pick session...")
tmux switch-client -t "$SESSION" || tmux attach -t "$SESSION"
```

- 从 `git` 历史中挑选一个提交哈希

```bash
git log --oneline | gum filter | cut -d' ' -f1 # | copy
```

- 简易的 [`skate`](https://github.com/charmbracelet/skate) 密码选择器

```
skate list -k | gum filter | xargs skate get
```

- 卸载软件包

```bash
brew list | gum choose --no-limit | xargs brew uninstall
```

- 清理 `git` 分支

```bash
git branch | cut -c 3- | gum choose --no-limit | xargs git branch -D
```

- 用 [`gh`](https://cli.github.com/) 检出 GitHub 上的 Pull Request

```bash
gh pr list | cut -f1,2 | gum choose | cut -f1 | xargs gh pr checkout
```

- 从 shell 历史中复制某条命令

```bash
gum filter < $HISTFILE --height 20
```

- `sudo` 的替代品

```bash
alias please="gum input --password | sudo -nS"
```

## 参与贡献

参见 [contributing][contribute]。

[contribute]: https://github.com/charmbracelet/gum/contribute

## 反馈

我们非常想听听你对这个项目的想法。欢迎随时给我们留言！

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## 许可证

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

---

[Charm](https://charm.sh) 的一部分。

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
