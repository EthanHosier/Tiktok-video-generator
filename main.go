package main

import "github.com/ethanhosier/clips/youtube"

const id = "Tb8WYOB89i4"

func main() {
	yt := youtube.NewYoutubeClient()
	yt.VideoForId(id)
}
