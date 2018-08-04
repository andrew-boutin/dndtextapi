// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/characters"
	log "github.com/sirupsen/logrus"

	"github.com/andrew-boutin/dndtextapi/channels"

	"github.com/andrew-boutin/dndtextapi/messages"
	"github.com/gin-gonic/gin"
)

// RegisterMessagesRoutes registers all of the Message routes with their
// associated middleware.
func RegisterMessagesRoutes(g *gin.RouterGroup) {
	g.GET("/channels/:channelID/messages", ValidateHeaders(acceptHeader), LoadChannelFromPathID, GetMessages)
	g.POST("/channels/:channelID/messages", ValidateHeaders(acceptHeader, contentTypeHeader), LoadChannelFromPathID, CreateMessage)
	g.GET("/channels/:channelID/messages/:id", ValidateHeaders(acceptHeader), LoadChannelFromPathID, GetMessage)
	g.PUT("/channels/:channelID/messages/:id", ValidateHeaders(acceptHeader, contentTypeHeader), UpdateMessage)
	g.DELETE("/channels/:channelID/messages/:id", DeleteMessage)
}

// GetMessages retrieves a list of Messages from the designated Channel. The query
// parameter msgType is optional and can be used to filter which Messages are
// retrieved.
func GetMessages(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)

	msgType, err := QueryParamExtractor(c, msgTypeQueryParam)
	if err != nil {
		// Query parameter is optional here so ignore not found error
		if err != ErrQueryParamNotFound {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	// Private Channels require that the User be a member to access any Messages
	// Accessing any meta Messages on public Channels also requires membership
	var isMember bool
	if channel.IsPrivate || msgType != storyMsgType {
		isMember, err = dbBackend.DoesUserHaveCharacterInChannel(user.ID, channel.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User is not a member of the Channel so deny access
		if !isMember {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}

	var onlyStoryMsgs bool
	var outMessages messages.MessageCollection
	switch msgType {
	case storyMsgType:
		onlyStoryMsgs = true
		outMessages, err = dbBackend.GetMessagesInChannel(channel.ID, &onlyStoryMsgs)
	case metaMsgType:
		onlyStoryMsgs = false
		outMessages, err = dbBackend.GetMessagesInChannel(channel.ID, &onlyStoryMsgs)
	default:
		outMessages, err = dbBackend.GetMessagesInChannel(channel.ID, nil)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, outMessages)
}

// GetMessage retrieves a single Message using the Message ID
// in the path.
func GetMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)

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

	// Private Channels and meta Messages in public Channels require
	// that the User is a member
	var isMember bool
	if channel.IsPrivate || !message.IsStory {
		isMember, err = dbBackend.DoesUserHaveCharacterInChannel(user.ID, channel.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// User is not a member of the Channel so deny access
		if !isMember {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}

	c.JSON(http.StatusOK, message)
}

// CreateMessage creates a new Message using the data in the
// request body.
func CreateMessage(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	channel := c.MustGet(channelKey).(*channels.Channel)

	// TODO: Validation - content not empty, can't set user/channel ids, etc.
	dbBackend := GetDBBackend(c)
	message := &messages.Message{}
	err := c.Bind(message)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// User must own the Character that the Message is for
	char, err := dbBackend.GetCharacter(message.CharacterID)
	if err != nil {
		if err == characters.ErrCharacterNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		log.WithError(err).Error("Failed to look up message.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.ID != char.UserID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Character must be in the Channel the Message is for
	if char.ChannelID != channel.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	message.ChannelID = channel.ID

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
	channel := c.MustGet(channelKey).(*channels.Channel)

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
		log.WithError(err).Error("Failed to look up message.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Look up the Character the Message is from so we can check if this User owns it
	char, err := dbBackend.GetCharacter(message.CharacterID)
	if err != nil {
		if err == characters.ErrCharacterNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		log.WithError(err).Error("Failed to look up character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// User must either have created the Message or be the Channel owner
	// to delete the Message
	if char.UserID != user.ID {
		// User didn't create the Message and isn't the Channel owner so deny access
		if channel.OwnerID != user.ID {
			c.AbortWithStatus(http.StatusForbidden)
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
		log.WithError(err).Error("Failed to look up message.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Look up the Character the Message is from so we can check if this User owns it
	char, err := dbBackend.GetCharacter(existingMessage.CharacterID)
	if err != nil {
		if err == characters.ErrCharacterNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		log.WithError(err).Error("Failed to look up character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// The User must have created the Message in order to update it
	if char.UserID != user.ID { // TODO: Check if User owns the Character tied to the Message *****
		c.AbortWithStatus(http.StatusForbidden)
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
