// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/andrew-boutin/dndtextapi/channels"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// RegisterChannelsRoutes registers all of the channel routes with their
// associated middleware
func RegisterChannelsRoutes(g *gin.RouterGroup) {
	g.GET("/channels", RequiredHeadersMiddleware(acceptHeader), GetChannels)
	g.POST("/channels", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), CreateChannel)
	g.GET("/channels/:channelID", RequiredHeadersMiddleware(acceptHeader), GetChannel)
	g.PUT("/channels/:channelID", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateChannel)
	g.DELETE("/channels/:channelID", DeleteChannel)
}

// GetChannels retrieves a list of Channels that the authenticated User has access
// to. An optional query parameter `level` can be used to filter Channels by if
// the User is a `member` or `owner`.
// TODO: What about query param visibility=public|private?
func GetChannels(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	// Check for the optional query parameter
	level, err := QueryParamExtractor(c, levelQueryParam)
	if err != nil {
		// Query parameter is optional here so ignore not found error
		if err != ErrQueryParamNotFound {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Get the Channels dependent on the query parameter
	var outChannels channels.ChannelCollection
	switch level {
	case ownerLevel:
		outChannels, err = dbBackend.GetChannelsOwnedByUser(user.ID)
	case memberLevel:
		outChannels, err = GetChannelsUserIsMember(dbBackend, user.ID)
	default:
		outChannels, err = GetChannelsUserCanAccess(dbBackend, user.ID)
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, outChannels)
}

// GetChannel retrieves a single Channel by using an id in the request path.
func GetChannel(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, channelIDPathParam)
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
	var userInChannel bool
	if channel.IsPrivate == true {
		userInChannel, err = dbBackend.DoesUserHaveCharacterInChannel(user.ID, channelID)
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
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channel := &channels.Channel{}
	err := c.Bind(channel)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Set the authenticated User as the Channel owner
	channel.OwnerID = user.ID

	createdChannel, err := dbBackend.CreateChannel(channel, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdChannel)
}

// DeleteChannel deletes the channel using the id from the request path.
func DeleteChannel(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, channelIDPathParam)
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
	if existingChannel.OwnerID != user.ID {
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
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, channelIDPathParam)
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
	if existingChannel.OwnerID != user.ID {
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

// GetChannelsUserIsMember finds all of the Channels that the User is a member of which
// means any Channel that the User has a Character in or owns.
func GetChannelsUserIsMember(dbBackend backends.Backend, userID int) (channels.ChannelCollection, error) {
	// User is considered a member of any Channel that they own
	ownedChannels, err := dbBackend.GetChannelsOwnedByUser(userID)
	if err != nil {
		log.WithError(err).Error("Failed to look up channels owned by user.")
		return nil, err
	}

	// Find all Channels that the User has a Character in
	charChannels, err := dbBackend.GetChannelsUserHasCharacterIn(userID, nil)
	if err != nil {
		log.WithError(err).Error("Failed to look up channels that user has a character in.")
		return nil, err
	}

	return concatUniqueChannels(ownedChannels, charChannels), nil
}

// GetChannelsUserCanAccess finds all Channels that a User has access to. This includes all public Channels,
// any private Channels that they have a Character in, and also any Channels that they own.
func GetChannelsUserCanAccess(dbBackend backends.Backend, userID int) (channels.ChannelCollection, error) {
	// Look up all public channels
	isPrivate := false
	publicChannels, err := dbBackend.GetAllChannels(&isPrivate)
	if err != nil {
		log.WithError(err).Error("Failed to look up public channels.")
		return nil, err
	}

	// Look up private Channels that the User has a Character in
	isPrivate = true
	charChannels, err := dbBackend.GetChannelsUserHasCharacterIn(userID, &isPrivate)
	if err != nil {
		log.WithError(err).Error("Failed to look up channels that user has a character in.")
		return nil, err
	}

	// Should be no intersection between public Channels and private Channels where there is a Character
	outChannels := append(publicChannels, charChannels...)

	// Look up channels owned by User
	ownedChannels, err := dbBackend.GetChannelsOwnedByUser(userID)
	if err != nil {
		log.WithError(err).Error("Failed to look up channels owned by user.")
		return nil, err
	}

	// Could be overlap between the Channels found previously and the Channels that are owned
	outChannels = concatUniqueChannels(outChannels, ownedChannels)
	return outChannels, nil
}

func concatUniqueChannels(cca channels.ChannelCollection, ccb channels.ChannelCollection) (newCC channels.ChannelCollection) {
	m := make(map[int]*channels.Channel)
	tmpCC := append(cca, ccb...)

	for _, c := range tmpCC {
		if _, ok := m[c.ID]; !ok {
			m[c.ID] = c
		}
	}

	for _, v := range m {
		newCC = append(newCC, v)
	}

	return newCC
}
