/*
 * MumbleDJ fork
 * By Richard Nysäter
 * commands/greet.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"
	"github.com/RichardNysater/mumbledj/services"
	"os"
)

var SUPPORTED_AUDIO_SUFFIXES = [...]string{"", ".mp3", ".ogg", ".wav", ".m4a", ".webm", ".opus", ".flac", ".aiff", ".aac", ".avi"}

// GreetCommand is a command that plays an audio track as a greeting to a user
type GreetCommand struct{}

// Aliases returns the current aliases for the command.
func (c *GreetCommand) Aliases() []string {
	return viper.GetStringSlice("commands.greet.aliases")
}

// Description returns the description for the command.
func (c *GreetCommand) Description() string {
	return viper.GetString("commands.greet.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *GreetCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.greet.is_admin")
}

// Execute executes the command with the given user and arguments.
// Return value descriptions:
//    string: A message to be returned to the user upon successful execution.
//    bool:   Whether the message should be private or not. true = private,
//            false = public (sent to whole channel).
//    error:  An error message to be returned upon unsuccessful execution.
//            If no error has occurred, pass nil instead.
// Example return statement:
//    return "This is a private message!", true, nil
func (c *GreetCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	if len(args) == 0 {
		return "", true, errors.New(viper.GetString("commands.greet.messages.no_user_error"))
	}
	if DJ.AudioStream == nil {
		username := args[0]
		logrus.Infoln("Attempting to greet user...")
		filepath, err := getGreetingFilepath(username)
		if err == nil {
			localFile := services.NewLocalFileService()
			track, err := localFile.GetTrack(filepath, username)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err.Error(),
					"filepath": filepath,
				}).Errorln("Could not get local file to play for greeting")
				return "", false, errors.New(viper.GetString("commands.greet.default_greeting_missing"))
			} else {
				DJ.Queue.InsertTrack(0, track)
			}
		}
	} else {
		logrus.Infoln("Did not greet user since bot is already playing.")
	}

	return "", false, nil
}

// GetGreetingFilepath will attempt to return the filepath to the greeting which should be played
func getGreetingFilepath(username string) (string, error) {
	filepath := os.ExpandEnv(viper.GetString("greetings.directory") + "/" + username)
	for _, element := range SUPPORTED_AUDIO_SUFFIXES {
		if _, err := os.Stat(filepath + element); !os.IsNotExist(err) {
			filepath += element
			return filepath, nil
		}
	}
	logrus.WithFields(logrus.Fields{
		"filepath": filepath,
	}).Info("No personal greeting audio file found at filepath. Playing default greeting instead.")

	filepath = os.ExpandEnv(viper.GetString("greetings.directory") + "/" +
		viper.GetString("greetings.default_filename"))

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"filepath": filepath,
		}).Errorln("Greeting is enabled but no default greeting audio file was found.")
		return "", err
	}
	return filepath, nil
}
