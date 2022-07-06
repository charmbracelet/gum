#!/bin/bash

go install github.com/charmbracelet/sodapop && mv $GOBIN/sodapop $GOBIN/pop

echo "Hello, there! Welcome to $(pop style --foreground 212 'Soda Pop')"

NAME=$(pop input --placeholder "What is your name?")

echo "Well, it is nice to meet you, $(pop style --foreground 212 "$NAME")."

COLOR=$(pop input --placeholder "What is your favorite color? (#HEX)")

echo "Wait a moment, while I think of my favorite color..."

pop spin --title "Thinking..." --color 212 -- sleep 3

echo "I like $(pop style --background $COLOR $COLOR), too. In fact, it's my $(pop style --background $COLOR 'favorite color!')"

sleep 1

echo "Seems like we have a lot in common, $(pop style --foreground 212 "$NAME")."

sleep 1

echo "What's your favorite Soda Pop flavor?"

POP=$(pop search --accent-color 212 << POPS
Cherry
Grape
Lime
Orange
POPS)

echo "One sec, while I finish my drink."

pop spin --title "Drinking some $POP soda pop..." --color 212 -- sleep 5

pop style --width 50 --padding "1 5" --margin "1 2" --border double --border-foreground 212 \
    "Well, it was nice meeting you, $(pop style --foreground 212 "$NAME"). Hope to see you soon!"\
    "Don't forget to drink some $(pop style --foreground 212 $POP) soda pop."
