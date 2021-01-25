package slack

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

// DefaultSlack keeps the Client API for Slack
type DefaultSlack struct {
	API *slack.Client
}

// NewDefaultSlack returns an access to Slack
func NewDefaultSlack() *DefaultSlack {
	godotenv.Load()
	token := os.Getenv("SLACK_TOKEN")
	return &DefaultSlack{
		API: slack.New(token),
	}
}

// GetChannelOrGroup returns a Channel ID for a named group or channel
func (s *DefaultSlack) GetChannelOrGroup(name string) (string, error) {
	if s == nil || s.API == nil {
		return "", errors.New("NotConnected")
	}
	next := ""
	for {
		conversation := &slack.GetConversationsParameters{
			Cursor:          next,
			ExcludeArchived: "true",
			Limit:           100,
			Types:           []string{"public_channel", "private_channel"},
		}
		channels, next, err := s.API.GetConversations(conversation)
		if err != nil {
			return "", nil
		}
		for _, v := range channels {
			if v.Name == name {
				return v.ID, nil
			}
		}
		if next == "" {
			return "", errors.New("NotFound")
		}
	}
}

func test() {
	s := NewDefaultSlack()
	channel, err := s.GetChannelOrGroup("mysql")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	s.API.PostMessage(
		channel, slack.MsgOptionText(
			"Notification",
			false,
		),
	)
	fmt.Printf("ID: %s\n", channel)
}
