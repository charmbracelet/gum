#!/bin/bash

gum style --border normal \
    --margin "1" \
    --padding "1 2" \
    --border-foreground 212 \
    "Hello, there! Welcome to $(gum style --foreground 212 'Gum')."

NAME=$(gum input --placeholder "What is your name?")

echo -e "Well, it is nice to meet you, $(gum style --foreground 212 "$NAME")."

sleep 2

clear

echo -e "Can you tell me a $(gum style --italic --foreground 99 'secret')?\n"

gum write --placeholder "I'll keep it to myself, I promise!" > /dev/null # we keep the secret to ourselves

clear

echo "What should I do with this information?"

sleep 1

READ="Read"
THINK="Think"
DISCARD="Think"
ACTIONS=$(gum choose --cursor-prefix "[ ] " --selected-prefix "[*] " --no-limit "$READ" "$THINK" "$DISCARD")

clear

echo "One moment, please."

if grep -q "$READ" <<< "$ACTION"; then
    gum spin -s line --title "Reading the secret..." -- sleep 1
fi

if grep -q "$THINK" <<< "$ACTION"; then
    gum spin -s pulse --title "Thinking about your secret..." -- sleep 1
fi

if grep -q "$DISCARD" <<< "$ACTION"; then
    gum spin -s monkey --title " Discarding your secret..." -- sleep 2
fi

sleep 1

clear

echo "What's your favorite $(gum style --foreground 212 "Gum") flavor?"

GUM=$(echo -e "Cherry\nGrape\nLime\nOrange" | gum filter)

echo "I'll keep that in mind!"

sleep 1

clear

echo "Do you like $(gum style --foreground "#04B575" "Bubble Gum?")"

sleep 1

CHOICE=$(gum choose --item.foreground 250 "Yes" "No" "It's complicated")

if [ "$CHOICE" == "Yes" ]; then
    echo "I thought so, $(gum style --bold "Bubble Gum") is the best."
else
    echo "I'm sorry to hear that."
fi

sleep 1

gum spin --title "Chewing some $(gum style --foreground "#04B575" "$GUM") bubble gum..." -- sleep 5

clear

gum join --horizontal \
    "$(gum style --height 5 --width 25 --padding '1 3' --border double --border-foreground 57 "Well, it was nice meeting you, $(gum style --foreground 212 "$NAME"). Hope to see you soon!")" \
    "$(gum style --width 25 --padding '1 3' --border double --border-foreground 212 "Don't forget to chew some $(gum style --foreground "#04B575" "$GUM") bubble gum.")"
