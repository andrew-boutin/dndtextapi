package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterAnonymousRoutes adds the anonymous routes
func RegisterAnonymousRoutes(r *gin.Engine) {
	g := r.Group("/public")
	g.GET("/channels", RequiredHeadersMiddleware(acceptHeader), GetPublicChannels)
	g.GET("/channels/:id", RequiredHeadersMiddleware(acceptHeader), GetPublicChannel)
	g.GET("/messages", RequiredHeadersMiddleware(acceptHeader), GetStoryMessagesInChannel)
}

// GetPublicChannels retrieves all of the public Channels accessible
// to anonymous Users.
func GetPublicChannels(c *gin.Context) {
	dbBackend := GetDBBackend(c)
	isPrivate := false
	channels, err := dbBackend.GetAllChannels(&isPrivate)
	if err != nil {
		log.WithError(err).Error("Failed to get public channels.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, channels)
}

// GetPublicChannel retrieves the Channel specified by the id in the
// path if it's public.
func GetPublicChannel(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	channel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			log.WithError(err).Error("Channel not found.")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to get channel.")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if channel.IsPrivate {
		log.WithError(err).Error("Anonymous user attempting to look up private channel denying access.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// GetStoryMessagesInChannel retrieves all of the story Messages from
// the Channel, if it's public, matching the id provided by the required
// query parameter channelID
func GetStoryMessagesInChannel(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := QueryParamAsIntExtractor(c, channelIDQueryParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingChannel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		// TODO: Maybe 400 if the channel id is bad
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if existingChannel.IsPrivate {
		log.WithError(err).Error("Anonymous User attempting to look up messages from private channel denying access.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	onlyStoryMsgs := true
	messages, err := dbBackend.GetMessagesInChannel(channelID, &onlyStoryMsgs)
	if err != nil {
		log.WithError(err).Error("Failed to get story messages for public channel.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, messages)
}