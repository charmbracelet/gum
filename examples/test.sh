#!/bin/sh

# Choose
gum choose Choose One Item --cursor "* " --cursor.foreground 99 --selected.foreground 99
gum choose Foo Bar Baz
gum choose Pick Two Items Maximum --limit 2 --cursor "* " --cursor-prefix "(•) " --selected-prefix "(x) " --unselected-prefix "( ) " --cursor.foreground 99 --selected.foreground 99
gum choose Unlimited Choice Of Items --no-limit --cursor "* " --cursor-prefix "(•) " --selected-prefix "(x) " --unselected-prefix "( ) " --cursor.foreground 99 --selected.foreground 99

# Confirm
gum confirm "Testing?"
gum confirm "No?" --default=false --affirmative "Okay." --negative "Cancel."

# Filter
echo -e {1..500} | sed 's/ /\n/g' | gum filter
echo -e {1..500} | sed 's/ /\n/g' | gum filter --indicator ">" --placeholder "Pick a number..." --indicator.foreground 1 --text.foreground 2 --match.foreground 3 --prompt.foreground 4 --height 5

# Format
echo -e "# Header\nBody" | gum format 
echo -e 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, Gum!")\n}\n' | gum format -t code
echo -e ":candy:" | gum format -t emoji
echo -e '{{ Bold "Bold" }}' | gum format -t template

# Input
gum input --prompt "Email: " --placeholder "john@doe.com" --prompt.foreground 99 --cursor.foreground 99 --width 50
gum input --password --prompt "Password: " --placeholder "hunter2" --prompt.foreground 99 --cursor.foreground 99 --width 50

# Join
gum join "Horizontal" "Join"
gum join --vertical "Vertical" "Join"

# Spin
gum spin --spinner minidot --title "Loading..." --title.foreground 99 -- sleep 1
gum spin --show-output --spinner monkey --title "Loading..." --title.foreground 99 -- sh -c 'sleep 1; echo "Hello, Gum!"'

# Style
gum style --foreground 99 --border double --border-foreground 99 --padding "1 2" --margin 1 "Hello, Gum."

# Write
gum write --width 40 --height 3 --placeholder "Type whatever you want" --prompt "| " --show-cursor-line --show-line-numbers --value "Something..." --base.padding 1 --cursor.foreground 99 --prompt.foreground 99
