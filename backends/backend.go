package backends

import (
	"log"

	"github.com/andrew-boutin/dndtextapi/backends/postgresql"
	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/configs"
	// Postgresql driver
	_ "github.com/lib/pq"
)

type Backend interface {
	// Channels functionality
	GetChannels() (*channels.ChannelCollection, error)
	GetChannel(int) (*channels.Channel, error)
	CreateChannel(*channels.Channel) (*channels.Channel, error)
}

func InitBackend(backendConfig configs.BackendConfiguration) (backendDB Backend) {
	log.Printf("backend config type is %s", backendConfig.Type)

	switch backendConfig.Type {
	case "postgres":
		backendDB = postgresql.GetPostgresqlBackend(backendConfig.User, backendConfig.PW, backendConfig.DBName)
	default:
		log.Fatalf("Unexpected backend config type %s", backendConfig.Type)
	}
	return
}
