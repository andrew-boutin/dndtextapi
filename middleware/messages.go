// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/andrew-boutin/dndtextapi/messages"
	"github.com/gin-gonic/gin"
)

// RegisterMessagesRoutes registers all of the Message routes with their
// associated middleware.
func RegisterMessagesRoutes(g *gin.RouterGroup) {
	g.GET("/messages", RequiredHeadersMiddleware(acceptHeader), GetMessages)
	g.POST("/messages", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), CreateMessage)
	g.GET("/messages/:id", RequiredHeadersMiddleware(acceptHeader), GetMessage)
	g.PUT("/messages/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateMessage)
	g.DELETE("/messages/:id", DeleteMessage)
}

// GetMessages retrieves a list of Messages. The query parameter channelID is
// required to determine which Channel to get the Messages from. The query
// parameter msgType is optional and can be used to filter which Messages are
// retrieved.
func GetMessages(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channelID, err := QueryParamAsIntExtractor(c, channelIDQueryParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	msgType, err := QueryParamExtractor(c, msgTypeQueryParam)
	if err != nil {
		// Query parameter is optional here so ignore not found error
		if err != ErrQueryParamNotFound {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	existingChannel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		// TODO: Maybe 400 if the channel id is bad
		// TODO: What about 404
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Private Channels require that the User be a member to access any Messages
	// Accessing any meta Messages on public Channels also requires membership
	var isMember bool
	if existingChannel.IsPrivate || msgType != storyMsgType {
		isMember, err = dbBackend.IsUserInChannel(user.ID, channelID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User is not a member of the Channel so deny access
		if !isMember {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	var onlyStoryMsgs bool
	var messages *messages.MessageCollection
	switch msgType {
	case storyMsgType:
		onlyStoryMsgs = true
		messages, err = dbBackend.GetMessagesInChannel(channelID, &onlyStoryMsgs)
	case metaMsgType:
		onlyStoryMsgs = false
		messages, err = dbBackend.GetMessagesInChannel(channelID, &onlyStoryMsgs)
	default:
		messages, err = dbBackend.GetMessagesInChannel(channelID, nil)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetMessage retrieves a single Message using the Message ID
// in the path.
func GetMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	message, err := dbBackend.GetMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	channel, err := dbBackend.GetChannel(message.ChannelID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Private Channels and meta Messages in public Channels require
	// that the User is a member
	var isMember bool
	if channel.IsPrivate || !message.IsStory {
		isMember, err = dbBackend.IsUserInChannel(user.ID, channel.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User is not a member of the Channel so deny access
		if !isMember {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	c.JSON(http.StatusOK, message)
}

// CreateMessage creates a new Message using the data in the
// request body.
func CreateMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	// TODO: Validation - content not empty, etc.
	dbBackend := GetDBBackend(c)
	message := &messages.Message{}
	err := c.Bind(message)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Verify that the User is a member of the Channel
	isMember, err := dbBackend.IsUserInChannel(user.ID, message.ChannelID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !isMember {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	message.UserID = user.ID

	createdMessage, err := dbBackend.CreateMessage(message)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdMessage)
}

// DeleteMessage deletes the message matching the ID in the path.
func DeleteMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	message, err := dbBackend.GetMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// User must either have created the Message or be the Channel owner
	// to delete the Message
	var channel *channels.Channel
	if message.UserID != user.ID {
		channel, err = dbBackend.GetChannel(message.ChannelID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User didn't create the Message and isn't the Channel owner so deny access
		if channel.OwnerID != user.ID {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	err = dbBackend.DeleteMessage(messageID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateMessage updates the Message using the ID from the path with
// the data from the request body.
func UpdateMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingMessage, err := dbBackend.GetMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// The User must have created the Message in order to update it
	if existingMessage.UserID != user.ID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	message := &messages.Message{}
	err = c.Bind(message)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedMessage, err := dbBackend.UpdateMessage(messageID, message)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedMessage)
}
