package backends

import (
	"log"

	"github.com/andrew-boutin/dndtextapi/messages"

	"github.com/andrew-boutin/dndtextapi/backends/postgresql"
	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/configs"
	// Postgresql driver
	_ "github.com/lib/pq"
)

// Backend defines the functionality expected of a backend.
type Backend interface {
	// Channels functionality
	GetChannels() (*channels.ChannelCollection, error)
	GetChannel(int) (*channels.Channel, error)
	CreateChannel(*channels.Channel) (*channels.Channel, error)
	DeleteChannel(int) error
	UpdateChannel(int, *channels.Channel) (*channels.Channel, error)

	// Messages functionality
	GetMessagesInChannel(int) (*messages.MessageCollection, error)
	GetMessage(int) (*messages.Message, error)
	CreateMessage(*messages.Message) (*messages.Message, error)
	DeleteMessage(int) error
	UpdateMessage(int, *messages.Message) (*messages.Message, error)
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
