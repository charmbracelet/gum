#!/bin/bash

export LIST=$(cat <<END
Cow:Moo
Cat:Meow
Dog:Woof
END
)

ANIMAL=$(echo "$LIST" | cut -d':' -f1 | gum filter)
SOUND=$(echo "$LIST" | grep $ANIMAL | cut -d':' -f2)

echo "The $ANIMAL goes $SOUND"
