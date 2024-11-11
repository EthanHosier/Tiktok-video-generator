package ffmpeg

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	attentionclips "github.com/ethanhosier/clips/attention-clips"
)

type FfmpegHandler interface {
	RemoveAudio(inputFile, outputFile string) (string, error)
	ClipVideo(inputFile, outputFile, startTime, duration string) (string, error)
}

type FfmpegClient struct {
}

func NewFfmpegClient() *FfmpegClient {
	return &FfmpegClient{}
}

func (ff *FfmpegClient) RemoveAudio(inputFile, outputFile string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-an", "-c:v", "copy", outputFile)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile, nil
}

func (ff *FfmpegClient) ClipVideoRandomSecs(inputFile attentionclips.AttentionClipsVideo, clipLength int) (string, error) {
	duration, err := ff.getVideoDuration(string(inputFile))
	if err != nil {
		return "", fmt.Errorf("failed to get video duration: %v", err)
	}

	totalDurationSeconds, err := strconv.Atoi(strings.Split(duration, ".")[0]) // assuming duration is in seconds with decimal precision
	if err != nil {
		return "", fmt.Errorf("failed to parse video duration: %v", err)
	}

	if clipLength >= totalDurationSeconds {
		return "", errors.New("clip length must be shorter than video duration")
	}

	rand.Seed(time.Now().UnixNano())
	maxStart := totalDurationSeconds - clipLength
	randomStartSeconds := rand.Intn(maxStart)

	randomStartTime := secondsToTimeString(randomStartSeconds)
	clipDuration := secondsToTimeString(clipLength)

	return ff.ClipVideo(inputFile, randomStartTime, clipDuration)
}

func (ff *FfmpegClient) ClipVideo(inputFile attentionclips.AttentionClipsVideo, startTime, duration string) (string, error) {
	if !isValidTime(startTime) {
		return "", fmt.Errorf("invalid start time: %s", startTime)
	}

	if !isValidTime(duration) {
		return "", fmt.Errorf("invalid duration: %s", duration)
	}

	outputPath := outputPathFor(string(inputFile), MP4)
	fmt.Println("outputPath: ", outputPath)

	cmd := exec.Command("ffmpeg", "-ss", startTime, "-i", string(inputFile), "-t", duration, "-c", "copy", outputPath)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputPath, nil
}