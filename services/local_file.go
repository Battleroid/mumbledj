/*
 * MumbleDJ fork
 * By Richard Nysäter
 * bot/local_file.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"regexp"
	"strings"
	"layeh.com/gumble/gumble"
	"github.com/Sirupsen/logrus"
	"github.com/RichardNysater/mumbledj/bot"
	"github.com/RichardNysater/mumbledj/interfaces"
	"sync"
	"os"
	"os/exec"
	"time"
	"strconv"
	"math"
)

// LocalFile is a custom service for getting local files.
// NOTE: Due to potential security concerns users can't currently use this service to play tracks.
type LocalFile struct {
	*GenericService
}

// NewLocalFileService returns an initialized LocalFile service object.
func NewLocalFileService() *LocalFile {
	return &LocalFile{
		&GenericService{
			ReadableName: "LocalFile",
			Format:       "bestaudio",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`(?i)^localfile:\.+`),
			},
			PlaylistRegex: nil,
		},
	}
}

// No API key check required to get local files
func (lf *GenericService) CheckAPIKey() error {
	return nil
}

// GetTracks uses the passed "url" (e.g. localfile:path/to/file.mp3) to find the track associated with the path.
// An error is returned if any error occurs when getting the local file.
func (lf *LocalFile) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		filepath string
		err      error
		track    bot.Track
		tracks   []interfaces.Track
	)

	filepath = strings.SplitN(url, ":", 1)[1]

	track, err = lf.GetTrack(filepath, submitter.Name)
	if err != nil {
		return nil, err
	}
	tracks = append(tracks, track)
	return tracks, nil
}

// GetTrack creates a track from a filepath to a valid audio file
func (lf *LocalFile) GetTrack(filepath string, username string) (bot.Track, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return bot.Track{}, err
	}

	fileInfo, _ := os.Stat(filepath)
	filename := fileInfo.Name()

	duration, err := getDuration(filepath)
	if err != nil {
		return bot.Track{}, err
	}
	logrus.WithFields(logrus.Fields{
		"filename": filename,
		"duration": duration,
	}).Infoln("Found track...")
	return bot.Track{
		ID:             filename,
		URL:            filepath,
		Title:          "LocalFile",
		Author:         username,
		Submitter:      username,
		Service:        lf.ReadableName,
		Filename:       filename,
		ThumbnailURL:   "",
		Duration:       duration,
		PlaybackOffset: 0,
		Playlist:       nil,
		WaitGroup:      &sync.WaitGroup{},
	}, nil
}

// Get the duration using ffprobe which the user should have installed
func getDuration(filepath string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filepath)
	cmd.Stderr = os.Stderr
	duration, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return durationFromBytesOfFloat64(string(duration))
}

// Return a time.Duration type from a string representation of a float
func durationFromBytesOfFloat64(s string) (time.Duration, error) {
	secondsFloat,err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil{
		return 0, err
	}
	secondsString := strconv.Itoa(int(math.Ceil(secondsFloat)))+"s"

	logrus.WithFields(logrus.Fields{
		"duration": secondsString,
	}).Infoln("Calculated duration of track...")
	return time.ParseDuration(secondsString)
}
