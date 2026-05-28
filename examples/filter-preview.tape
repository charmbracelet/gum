Output filter-preview.gif
Set Width 1200
Set Height 600
Set Shell bash

Type "seq 1 10 | ./gum filter --height 12 --preview 'echo {} times 2 = $(({} * 2)); echo {} squared = $(({} * {}))'"
Enter
Sleep 500ms
Down@200ms 3
Sleep 500ms
Up@200ms 1
Sleep 500ms
Enter
Sleep 1s
