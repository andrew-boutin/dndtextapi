// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package backends

import (
	log "github.com/sirupsen/logrus"

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
}

// InitBackend initializes whatever backend matches the provided
// configuration and returns it.
func InitBackend(backendConfig configs.BackendConfiguration) (backendDB Backend) {
	log.Printf("backend config type is %s", backendConfig.Type)

	switch backendConfig.Type {
	case "postgres":
		backendDB = postgresql.MakePostgresqlBackend(backendConfig.User, backendConfig.PW, backendConfig.DBName)
	default:
		log.Fatalf("Unexpected backend config type %s", backendConfig.Type)
	}
	return
}
