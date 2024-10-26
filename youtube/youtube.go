package youtube

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/kkdai/youtube/v2"
)

type YoutubeHandler interface {
	CaptionsForId(id string)
}

type YoutubeClient struct {
}

func NewYoutubeClient() *YoutubeClient {
	return &YoutubeClient{}
}

func (y *YoutubeClient) CaptionsForId(id string) {
	serverErrCh := make(chan error)
	go startServer(serverErrCh, id)

	requestUrlCh := make(chan *url.URL)
	go simulateBrowserForUrl(requestUrlCh)

	// Handle any server errors from the serverErrCh channel
	select {
	case err := <-serverErrCh:
		log.Fatalf("Server error %v", err)
	case requestUrl := <-requestUrlCh:
		fmt.Println("Received request URL:", requestUrl)
	case <-time.After(10 * time.Second):
		log.Fatal("Timeout")
	}
}

func (y *YoutubeClient) VideoForId(id string) {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		log.Fatalf("Error getting video: %v", err)
	}

	fmt.Printf("%+v", video)
}

func handler(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<iframe width="560" height="315" src="https://www.youtube.com/embed/%s?si=jGSgSFl5DP3Rxyin&cc_load_policy=1" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
	</body>
	</html>`, id)
		fmt.Fprint(w, html)
	}
}

func startServer(errCh chan error, id string) {
	http.HandleFunc("/", handler(id))
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		errCh <- err
	}
}

func simulateBrowserForUrl(requestUrlCh chan *url.URL) {
	url := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	router := browser.HijackRequests()
	defer router.MustStop()

	router.MustAdd("*youtube.com/api/timedtext*", func(ctx *rod.Hijack) {
		requestUrlCh <- ctx.Request.URL()
		ctx.ContinueRequest(nil)
	})

	go router.Run()

	page := browser.MustPage("http://localhost:8080")
	page.MustWaitLoad()
	page.MustWaitRequestIdle()

	time.Sleep(5 * time.Second)
	page.Mouse.MustMoveTo(300, 300).MustClick("left")
	fmt.Println("Clicked on the video")
	time.Sleep(5 * time.Second)
}
