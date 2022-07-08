#!/bin/bash

go install github.com/charmbracelet/gum

echo "Hello, there! Welcome to $(gum style --foreground 212 'Gum')."

NAME=$(gum input --placeholder "What is your name?")

echo "Well, it is nice to meet you, $(gum style --foreground 212 "$NAME")."

COLOR=$(gum input --placeholder "What is your favorite color? (#HEX)")

echo "Wait a moment, while I think of my favorite color..."

gum spin --title "Thinking..." --color 212 -- sleep 3

echo "I like $(gum style --padding "0 1" --background $COLOR $COLOR), too. In fact, it's my $(gum style --padding "0 1" --background $COLOR 'favorite color!')"

sleep 1

echo "Seems like we have a lot in common, $(gum style --foreground 212 "$NAME")."

sleep 1

echo "What's your favorite Gum flavor?"

GUM=$(echo "Cherry
Grape
Lime
Orange" | gum filter)

echo "One sec, while I finish my drink."

gum spin --title "Chewing some $GUM bubble gum..." --color 212 -- sleep 5

gum style --width 50 --padding "1 5" --margin "1 2" --border double --border-foreground 212 \
    "Well, it was nice meeting you, $(gum style --foreground 212 "$NAME"). Hope to see you soon!"\
    "Don't forget to chew some $(gum style --foreground 212 $GUM) bubble gum."
