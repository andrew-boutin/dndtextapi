// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package backends

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/andrew-boutin/dndtextapi/bots"
	"github.com/andrew-boutin/dndtextapi/characters"
	"github.com/andrew-boutin/dndtextapi/messages"
	"github.com/andrew-boutin/dndtextapi/users"

	"github.com/andrew-boutin/dndtextapi/backends/postgresql"
	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/configs"
	// Postgresql driver
	_ "github.com/lib/pq"
)

// Backend defines the functionality expected of a backend.
type Backend interface {
	// Channels functionality
	GetChannel(int) (*channels.Channel, error)
	GetChannelsOwnedByUser(int) (channels.ChannelCollection, error)
	GetChannelsUserHasCharacterIn(int, *bool) (channels.ChannelCollection, error)
	GetAllChannels(*bool) (channels.ChannelCollection, error)
	CreateChannel(*channels.Channel, int) (*channels.Channel, error)
	DeleteChannel(int) error
	UpdateChannel(int, *channels.Channel) (*channels.Channel, error)

	// Messages functionality
	GetMessagesInChannel(int, *bool) (messages.MessageCollection, error)
	GetMessage(int) (*messages.Message, error)
	CreateMessage(*messages.Message) (*messages.Message, error)
	DeleteMessage(int) error
	UpdateMessage(int, *messages.Message) (*messages.Message, error)
	DeleteMessagesFromUser(int) error
	DeleteMessagesFromChannel(int) error
	DeleteMessagesFromCharacter(int) error

	// Users functionality
	UpdateUser(int, *users.User) (*users.User, error)
	DeleteUser(int) error
	GetUserByEmail(string) (*users.User, error)
	GetUserByID(int) (*users.User, error)
	CreateUser(*users.GoogleUser) (*users.User, error)
	GetAllUsers() (users.UserCollection, error)
	UpdateUserLastLogin(*users.User) (*users.User, error)

	// Characters functionality
	DoesUserHaveCharacterInChannel(int, int) (bool, error)
	GetCharactersInChannel(channelID int) (characters.CharacterCollection, error)
	GetCharacter(int) (*characters.Character, error)
	CreateCharacter(*characters.Character) (*characters.Character, error)
	DeleteCharacter(int) error
	UpdateCharacter(int, *characters.Character) (*characters.Character, error)
	DeleteCharactersFromUser(int) error
	DeleteCharactersFromChannel(int) error

	// Bots functionality
	GetAllBots() (bots.BotCollection, error)
	GetBot(int) (*bots.Bot, error)
	GetBotCreds(int) (*bots.BotClientCredentials, error)
	DeleteBot(int) error
	CreateBot(*bots.Bot) (*bots.Bot, error)
	CreateBotCreds(*bots.Bot) (*bots.BotClientCredentials, error)
	UpdateBot(int, *bots.Bot) (*bots.Bot, error)
	DeleteBotCreds(int) error
}

// InitBackend initializes whatever backend matches the provided
// configuration and returns it.
func InitBackend(backendConfig configs.BackendConfiguration) (backendDB Backend, err error) {
	switch backendConfig.Type {
	case "postgres":
		backendDB, err = postgresql.MakePostgresqlBackend(backendConfig.User, backendConfig.PW, backendConfig.DBName)
		if err != nil {
			log.WithError(err).Error("Failed to initialize postgresql backend.")
		}
	default:
		err = fmt.Errorf("Unexpected backend config type %s", backendConfig.Type)
		log.WithError(err).Error("Failed to initialize a backend.")
	}
	return
}
