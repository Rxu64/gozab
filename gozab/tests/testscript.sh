#!  /usr/bin/osascript

tell application "Terminal"
    set window9 to do script "cd gozab/gozab"
    delay 0.2
    do script "tests/testLaunchService.sh" in window9
end tell

delay 2

tell application "Terminal"
    set window10 to do script "cd gozab/gozab"
    delay 0.2
    do script "tests/testUserA.sh" in window10
end tell