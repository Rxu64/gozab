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
    do script "go run node/main.go localhost:50051 localhost:50056" in window1
    do script "go run node/main.go localhost:50052 localhost:50057" in window2
    do script "go run node/main.go localhost:50053 localhost:50058" in window3
    do script "go run node/main.go localhost:50054 localhost:50059" in window4
    do script "go run node/main.go localhost:50055 localhost:50060" in window5
end tell