#!/bin/bash

# This script converts some video to a GIF. It prompts the user to select an
# video file with `gum filter` Set the frame rate, desired width, and max
# colors to use Then, converts the video to a GIF.

INPUT=$(gum filter --placeholder "Input file")
FRAMERATE=$(gum input --prompt "Frame rate: " --placeholder "Frame Rate" --prompt.foreground 240 --value "50")
WIDTH=$(gum input --prompt "Width: " --placeholder "Width" --prompt.foreground 240 --value "1200")
MAXCOLORS=$(gum input --prompt "Max Colors: " --placeholder "Max Colors" --prompt.foreground 240 --value "256")

BASENAME=$(basename "$INPUT")
BASENAME="${BASENAME%%.*}"

gum spin --title "Converting to GIF" -- ffmpeg -i "$INPUT" -vf "fps=$FRAMERATE,scale=$WIDTH:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=$MAXCOLORS[p];[s1][p]paletteuse" "$BASENAME.gif"
