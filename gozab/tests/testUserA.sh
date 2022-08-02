#!  /usr/bin/osascript

tell application "Terminal"
    set window6 to do script "cd gozab/gozab"
    delay 0.3
    do script "go mod tidy" in window6
    do script "go run user/main.go" in window6
    do script "Send" in window6
    do script "a 1" in window6
    do script "Send" in window6
    do script "b 2" in window6
    do script "Send" in window6
    do script "c 3" in window6
end tell

delay 1

tell application "Terminal"
    set window7 to do script "cd gozab/gozab"
    delay 0.3
    do script "go mod tidy" in window7
    do script "go run user/main.go" in window7
    do script "Get" in window7
    do script "a" in window7
    do script "Get" in window7
    do script "b" in window7
    do script "Get" in window7
    do script "c" in window7
end tell