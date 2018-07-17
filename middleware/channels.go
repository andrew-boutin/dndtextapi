package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/gin-gonic/gin"
)

// RegisterChannelsMiddleware registers all of the channel routes with their
// associated middleware
func RegisterChannelsMiddleware(r *gin.Engine) {
	r.GET("/channels", RequiredHeadersMiddleware(acceptHeader), GetChannels)
	r.POST("/channels", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), CreateChannel)
	r.GET("/channels/:id", RequiredHeadersMiddleware(acceptHeader), GetChannel)
	r.PUT("/channels/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateChannel)
	r.DELETE("/channels/:id", DeleteChannel)
}

// GetChannels retrieves a list of Channels that the authenticated User has access
// to. An optional query parameter `level` can be used to filter Channels by if
// the User is a `member` or `owner`.
// TODO: What about query param visibility=public|private?
func GetChannels(c *gin.Context) {
	userID := GetAuthenticatedUserID()
	dbBackend := GetDBBackend(c)

	// Check for the optional query parameter
	level, err := QueryParamExtractor(c, levelQueryParam)
	if err != nil {
		// Query parameter is optional here so ignore not found error
		if err != ErrQueryParamNotFound {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	// Get the Channels dependent on the query parameter
	var channels *channels.ChannelCollection
	switch level {
	case ownerLevel:
		channels, err = dbBackend.GetChannelsOwnedByUser(userID)
	case memberLevel:
		channels, err = dbBackend.GetChannelsUserIsMember(userID, nil)
	default:
		isPrivate := false
		publicChannels, err := dbBackend.GetAllChannels(&isPrivate)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		isPrivate = true
		privateMemberChannels, err := dbBackend.GetChannelsUserIsMember(userID, &isPrivate)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		usersChannels := append(*publicChannels, *privateMemberChannels...)
		channels = &usersChannels
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, channels)
}

// GetChannel retrieves a single Channel by using an id in the request path.
func GetChannel(c *gin.Context) {
	userID := GetAuthenticatedUserID()
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
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

	// Private Channels require that the User is a member to access
	if channel.IsPrivate == true {
		userInChannel, err := dbBackend.IsUserInChannel(userID, channelID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User not a member of the Channel so deny access
		if !userInChannel {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	c.JSON(http.StatusOK, channel)
}

// CreateChannel creates a new channel using the data provided
// in the request body.
func CreateChannel(c *gin.Context) {
	// TODO: Validation. Name not empty, can't set ID/OwnerID, etc.
	userID := GetAuthenticatedUserID()
	dbBackend := GetDBBackend(c)

	channel := &channels.Channel{}
	err := c.Bind(channel)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Set the authenticated User as the Channel owner
	channel.OwnerID = userID

	createdChannel, err := dbBackend.CreateChannel(channel, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add the authenticated User as a Channel member
	err = dbBackend.AddUserToChannel(userID, createdChannel.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdChannel)
}

// DeleteChannel deletes the channel using the id from the request path.
func DeleteChannel(c *gin.Context) {
	userID := GetAuthenticatedUserID()
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingChannel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// User must be the owner of the Channel in order to delete
	if existingChannel.OwnerID != userID {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	err = dbBackend.DeleteChannel(channelID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateChannel updates the specified channel from the id in the request
// path using the data in the request body.
func UpdateChannel(c *gin.Context) {
	// TODO: Validation. Name not empty, can't set ID/OwnerID, etc.
	userID := GetAuthenticatedUserID()
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
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

	existingChannel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// User must be the owner of the Channel in order to update
	if existingChannel.OwnerID != userID {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	updatedChannel, err := dbBackend.UpdateChannel(channelID, channel)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedChannel)
}
