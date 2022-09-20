#!  /usr/bin/osascript

tell application "Terminal"
    set window1 to do script "cd gozab"
    set window2 to do script "cd gozab"
    set window3 to do script "cd gozab"
    set window4 to do script "cd gozab"
    set window5 to do script "cd gozab"
    delay 1
    do script "go mod tidy" in window1
    do script "go mod tidy" in window2
    do script "go mod tidy" in window3
    do script "go mod tidy" in window4
    do script "go mod tidy" in window5
    do script "go run node/main.go 0" in window1
    do script "go run node/main.go 1" in window2
    do script "go run node/main.go 2" in window3
    do script "go run node/main.go 3" in window4
    do script "go run node/main.go 4" in window5
end tell