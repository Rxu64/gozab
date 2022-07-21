package main

import (
	f "gozab/follower"
	"os"
)

func main() {
	go f.FollowerRoutine(os.Args)
	for {
		// I know this looks weird
		// this main need to be alive to keep the routine alive
	}
}
