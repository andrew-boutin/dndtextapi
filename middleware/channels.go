package middleware

import (
	"net/http"
	"strconv"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/gin-gonic/gin"
)

func RegisterChannelsMiddleware(r *gin.Engine) {
	r.GET("/channels", GetChannels)
	r.GET("/channels/:id", GetChannel)
	r.POST("/channels", CreateChannel)
}

func GetChannels(c *gin.Context) {
	// TODO: Anonymous users can get a list of public channels. Users who are authn
	// can also see private channels they're a member of.
	dbBackend := GetDBBackend(c)
	channels, err := dbBackend.GetChannels()
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, channels)
}

func GetChannel(c *gin.Context) {
	// TODO: Anonymous users can get a public channel. Users who are authn
	// can also get a private channel they're a member of.
	dbBackend := GetDBBackend(c)
	idParam := c.Param("id")
	channelID, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(err)
		return
	}
	channel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, channel)
}

// TODO: Validation. Name not empty etc.
func CreateChannel(c *gin.Context) {
	dbBackend := GetDBBackend(c)
	channel := &channels.Channel{}
	err := c.Bind(channel)

	if err != nil {
		c.Error(err)
		return
	}

	createdChannel, err := dbBackend.CreateChannel(channel)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdChannel)
}
