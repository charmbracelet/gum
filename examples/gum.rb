puts 'What is your name?'
name = `gum input --placeholder "Your name"`.chomp

puts "Hello #{name}!"

puts 'Pick your 2 favorite colors'

COLORS = {
  'Red' => '#FF0000',
  'Blue' => '#0000FF',
  'Green' => '#00FF00',
  'Yellow' => '#FFFF00',
  'Orange' => '#FFA500',
  'Purple' => '#800080',
  'Pink' => '#FF00FF'
}.freeze

colors = `gum choose #{COLORS.keys.join(' ')} --limit 2`.chomp.split("\n")

if colors.length == 2
  first = `gum style --foreground '#{COLORS[colors[0]]}' '#{colors[0]}'`.chomp
  second = `gum style --foreground '#{COLORS[colors[1]]}' '#{colors[1]}'`.chomp
  puts "You chose #{first} and #{second}."
elsif colors.length == 1
  first = `gum style --foreground '#{COLORS[colors[0]]}' '#{colors[0]}'`.chomp
  puts "You chose #{first}."
else
  puts "You didn't pick any colors!"
end
