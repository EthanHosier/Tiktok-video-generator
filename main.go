package main

import (
	"fmt"

	"github.com/ethanhosier/clips/youtube"
)

const id = "Tb8WYOB89i4"

func main() {
	yt := youtube.NewYoutubeClient()
	vid, err := yt.VideoForId(id)
	if err != nil {
		fmt.Printf("error getting video: %v\n", err)
		return
	}

	fmt.Printf("vid: %+v\n", vid)
}
