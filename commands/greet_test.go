/*
 * MumbleDJ fork
 * By Richard Nysäter
 * commands/greet_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"layeh.com/gumble/gumbleffmpeg"
	"github.com/RichardNysater/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type GreetCommandTestSuite struct {
	Command GreetCommand
	suite.Suite
}

func (suite *GreetCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()
	bot.DJ = DJ

	// Trick the tests into thinking audio is already playing to avoid
	// attempting to play tracks that don't exist.
	DJ.AudioStream = new(gumbleffmpeg.Stream)

	viper.Set("commands.greet.aliases", []string{"greet"})
	viper.Set("commands.greet.description", "greet")
	viper.Set("commands.greet.is_admin", false)
}

func (suite *GreetCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *GreetCommandTestSuite) TestAliases() {
	suite.Equal([]string{"greet"}, suite.Command.Aliases())
}

func (suite *GreetCommandTestSuite) TestDescription() {
	suite.Equal("greet", suite.Command.Description())
}

func (suite *GreetCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

func (suite *GreetCommandTestSuite) TestExecuteWithNoArgs() {
	message, isPrivateMessage, err := suite.Command.Execute(nil)

	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned for attempting to greet without providing a user.")
}

// TODO: Implement this test.
func (suite *GreetCommandTestSuite) TestExecuteWhenPersonalGreetingFound() {

}

// TODO: Implement this test.
func (suite *GreetCommandTestSuite) TestExecuteWhenNoPersonalGreetingFound() {

}

// TODO: Implement this test.
func (suite *GreetCommandTestSuite) TestExecuteWhenNoPersonalOrDefaultGreetingFound() {

}

// TODO: Implement this test.
func (suite *GreetCommandTestSuite) TestExecuteWhenNoValidSuffixFound() {

}

func TestGreetCommandTestSuite(t *testing.T) {
	suite.Run(t, new(GreetCommandTestSuite))
}
