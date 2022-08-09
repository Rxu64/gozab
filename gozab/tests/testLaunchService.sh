#!  /usr/bin/osascript

tell application "Terminal"
    set window1 to do script "cd gozab/gozab"
    set window2 to do script "cd gozab/gozab"
    set window3 to do script "cd gozab/gozab"
    set window4 to do script "cd gozab/gozab"
    set window5 to do script "cd gozab/gozab"
    delay 1
    do script "go mod tidy" in window1
    do script "go mod tidy" in window2
    do script "go mod tidy" in window3
    do script "go mod tidy" in window4
    do script "go mod tidy" in window5
    do script "go run node/main.go localhost:50051" in window1
    do script "go run node/main.go localhost:50052" in window2
    do script "go run node/main.go localhost:50053" in window3
    do script "go run node/main.go localhost:50054" in window4
    do script "go run node/main.go localhost:50055" in window5
end tell