Output example3.gif

Set TypingSpeed 20ms
Set FontFamily "Cartograph CF"
Set FontSize 15
Set Width 1280
Set Height 720
Set Theme "Desert"

# Go to parent dir
Type "cd .."
Enter

# Show the current colorscheme of the app
Type "go run cmd/goread/main.go -colors"
Enter
Sleep 5s

# Change the currnent pywal scheme
Type "wal -stni /home/adam/pictures/wallpapers/Sunny\ outside.webp > /dev/null"
Enter
Sleep 3s

# Show the new colorscheme of the app
Type "go run cmd/goread/main.go -pywal"
Enter
Sleep 5s

# Run the app
Type "go run cmd/goread/main.go"
Enter
Sleep 5s
