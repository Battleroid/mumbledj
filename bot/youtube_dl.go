/*
 * MumbleDJ
 * By Matthieu Grieger
 * bot/youtube_dl.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package bot

import (
	"errors"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/RichardNysater/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// YouTubeDL is a struct that gathers all methods related to the youtube-dl
// software.
// youtube-dl: https://rg3.github.io/youtube-dl/
type YouTubeDL struct{}

// Download downloads the audio associated with the incoming `track` object
// and stores it `track.Filename`.
func (yt *YouTubeDL) Download(t interfaces.Track) error {
	t.GetWaitGroup().Add(1)
	player := "--prefer-ffmpeg"
	if viper.GetString("defaults.player_command") == "avconv" {
		player = "--prefer-avconv"
	}

	filepath := os.ExpandEnv(viper.GetString("cache.directory") + "/" + t.GetFilename())

	// Determine which format to use.
	format := "bestaudio"
	for _, service := range DJ.AvailableServices {
		if service.GetReadableName() == t.GetService() {
			format = service.GetFormat()
		}
	}

	// Check to see if track is already downloaded.
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"url": t.GetURL(),
			"download-path": filepath,
			"format" : format,
			"player" : player,
		}).Infoln("Downloading track...")

		var cmd *exec.Cmd
		if t.GetService() == "Mixcloud" {
			logrus.Infoln("Downloading from Mixcloud...")
			cmd = exec.Command("youtube-dl", "--verbose", "--no-mtime", "--output", filepath, "--format", format, "--external-downloader", "aria2c", player, t.GetURL())
		} else {
			cmd = exec.Command("youtube-dl", "--verbose", "--no-mtime", "--output", filepath, "--format", format, player, t.GetURL())
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			args := ""
			for s := range cmd.Args {
				args += cmd.Args[s] + " "
			}
			logrus.Warnf("%s\n%s\nyoutube-dl: %s", args, string(output), err.Error())
			t.GetWaitGroup().Done()
			return errors.New("Track download failed")
		}
		if viper.GetBool("cache.enabled") {
			DJ.Cache.CheckDirectorySize()
		}
	}
	t.GetWaitGroup().Done()
	return nil
}

// Delete deletes the audio file associated with the incoming `track` object.
func (yt *YouTubeDL) Delete(t interfaces.Track) error {
	if !viper.GetBool("cache.enabled") {
		filePath := os.ExpandEnv(viper.GetString("cache.directory") + "/" + t.GetFilename())
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err == nil {
				return nil
			}
			return errors.New("An error occurred while deleting the audio file")
		}
	}
	return nil
}
