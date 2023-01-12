Output example2.gif

Set TypingSpeed 20ms
Set FontFamily "Cartograph CF"
Set FontSize 15
Set Width 1280
Set Height 720
Set Theme "Desert"

# Show the current colorscheme of the app
Type "../goread --test_colors"
Enter
Sleep 5s

# Change the currnent pywal scheme
Type "wal -stni /home/adam/pictures/wallpapers/Sunny\ outside.webp > /dev/null"
Enter
Sleep 3s

# Show the new colorscheme of the app
Type "../goread --get_colors ~/.cache/wal/colors.json"
Enter
Sleep 5s

# Run the app
Type "../goread"
Enter
Sleep 5s
