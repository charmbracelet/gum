#!/bin/bash

# Always ask for permission!
echo "Do you want to see a magic trick?"

YES="Yes, please!"
NO="No, thank you!"

CHOICE=$(gum choose "$YES" "$NO")

if [ "$CHOICE" != "$YES" ]; then
    echo "Alright, then. Have a nice day!"
    exit 1
fi


# Let the magic begin.
echo "Alright, then. Let's begin!"
gum style --foreground 212 "Pick a card, any card..."

CARD=$(gum choose "Ace (A)" "Two (2)" "Three (3)" "Four (4)" "Five (5)" "Six (6)" "Seven (7)" "Eight (8)" "Nine (9)" "Ten (10)" "Jack (J)" "Queen (Q)" "King (K)")
SUIT=$(gum choose "Hearts (♥)" "Diamonds (♦)" "Clubs (♣)" "Spades (♠)")

gum style --foreground 212 "You picked the $CARD of $SUIT."

SHORT_CARD=$(echo $CARD | cut -d' ' -f2 | tr -d '()')
SHORT_SUIT=$(echo $SUIT | cut -d' ' -f2 | tr -d '()')

TOP_LEFT=$(gum join --vertical "$SHORT_CARD" "$SHORT_SUIT")
BOTTOM_RIGHT=$(gum join --vertical "$SHORT_SUIT" "$SHORT_CARD")

TOP_LEFT=$(gum style --width 10 --height 5 --align left "$TOP_LEFT")
BOTTOM_RIGHT=$(gum style --width 10 --align right "$BOTTOM_RIGHT")

if [[ "$SHORT_SUIT" == "♥" || "$SHORT_SUIT" == "♦" ]]; then
    CARD_COLOR="1" # Red
else
    CARD_COLOR="7" # Black
fi

gum style --border rounded --padding "0 1" --margin 2 --border-foreground "$CARD_COLOR" --foreground "$CARD_COLOR" "$(gum join --vertical "$TOP_LEFT" "$BOTTOM_RIGHT")"

echo "Is this your card?"

gum choose "Omg, yes!" "Nope, sorry!"
