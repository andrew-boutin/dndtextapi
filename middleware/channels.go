package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/gin-gonic/gin"
)

// RegisterChannelsMiddleware registers all of the channel routes with their
// associated middleware
func RegisterChannelsMiddleware(r *gin.Engine) {
	// TODO: Maybe don't include any user info in these calls and have that be separate
	// under /channels/:channelID/users ...
	r.GET("/channels", RequiredHeadersMiddleware(acceptHeader), GetChannels)
	r.POST("/channels", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), CreateChannel)
	r.GET("/channels/:id", RequiredHeadersMiddleware(acceptHeader), GetChannel)
	r.PUT("/channels/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateChannel)
	r.DELETE("/channels/:id", DeleteChannel)
}

// GetChannels retrieves a list of all channels. The channels returned
// are partial views.
func GetChannels(c *gin.Context) {
	// TODO: Anonymous users can get a list of public channels. Users who are authn
	// can also see private channels they're a member of.
	dbBackend := GetDBBackend(c)
	channels, err := dbBackend.GetChannels()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, channels)
}

// GetChannel retrieves a single channel with all of its data by using
// an id in the request path.
func GetChannel(c *gin.Context) {
	// TODO: Anonymous users can get a public channel. Users who are authn
	// can also get a private channel they're a member of.
	dbBackend := GetDBBackend(c)
	channelID, err := PathParamAsIntExtractor(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	channel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// CreateChannel creates a new channel using the data provided
// in the request body.
func CreateChannel(c *gin.Context) {
	// TODO: User must be authn
	// TODO: Validation. Name not empty etc.
	dbBackend := GetDBBackend(c)
	channel := &channels.Channel{}
	err := c.Bind(channel)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createdChannel, err := dbBackend.CreateChannel(channel)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdChannel)
}

// DeleteChannel deletes the channel using the id from the request path.
func DeleteChannel(c *gin.Context) {
	// TODO: Authn user must be owner
	dbBackend := GetDBBackend(c)
	channelID, err := PathParamAsIntExtractor(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dbBackend.DeleteChannel(channelID)

	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateChannel updates the specified channel from the id in the request
// path using the data in the request body.
func UpdateChannel(c *gin.Context) {
	// TODO: Authn user must be owner
	dbBackend := GetDBBackend(c)
	channelID, err := PathParamAsIntExtractor(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	channel := &channels.Channel{}
	err = c.Bind(channel)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedChannel, err := dbBackend.UpdateChannel(channelID, channel)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedChannel)
}
