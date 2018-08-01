// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"errors"
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/characters"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterCharactersRoutes registers all of the character routes with their
// associated middleware
func RegisterCharactersRoutes(g *gin.RouterGroup) {
	g.GET("/channels/:channelID/characters", RequiredHeadersMiddleware(acceptHeader), LoadChannelFromPathID, GetCharacters)
	g.POST("/channels/:channelID/characters", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), LoadChannelFromPathID, CreateCharacter)
	g.GET("/channels/:channelID/characters/:id", RequiredHeadersMiddleware(acceptHeader), LoadChannelFromPathID, LoadCharacter, GetCharacter)
	g.PUT("/channels/:channelID/characters/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), LoadChannelFromPathID, UpdateCharacter)
	g.DELETE("/channels/:channelID/characters/:id", LoadChannelFromPathID, LoadCharacter, DeleteCharacter)
}

// GetCharacters retrieves all of the Characters in the Channel from the path. The
// User must have a Character in the Channel or be the owner of the Channelß.
func GetCharacters(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)

	// If the User isn't the Channel owner then they have to have a Character in the Channel
	if user.ID != channel.OwnerID {
		isUserInChannel, err := dbBackend.DoesUserHaveCharacterInChannel(user.ID, channel.ID)
		if err != nil {
			log.WithError(err).Error("Failed to look up if user is in channel.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !isUserInChannel {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}

	charactersInChannel, err := dbBackend.GetCharactersInChannel(channel.ID)
	if err != nil {
		log.WithError(err).Error("Failed to look up characters for channel.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, charactersInChannel)
}

// CreateCharacter allows the Channel owner to create a new Character. This
// is how Users are invited to a Channel.
func CreateCharacter(c *gin.Context) {
	// TODO: Can't fill in name
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)

	character := &characters.Character{}
	err := c.Bind(character)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if channel.OwnerID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	newCharacter, err := dbBackend.CreateCharacter(character)
	if err != nil {
		log.WithError(err).Error("Failed to create character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newCharacter)
}

// GetCharacter retrieves a single character using the id from the path. Anyone with a
// Character is allowed along with the Channel owner.
func GetCharacter(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)
	character := c.MustGet(characterKey).(*characters.Character)

	// If the User isn't the Channel owner then they have to have a Character in the Channel
	if user.ID != channel.OwnerID {
		isUserInChannel, err := dbBackend.DoesUserHaveCharacterInChannel(user.ID, channel.ID)
		if err != nil {
			log.WithError(err).Error("Failed to look up if user is in channel.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !isUserInChannel {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}

	c.JSON(http.StatusOK, character)
}

// UpdateCharacter allows the Character owner to update their Character. The Character
// is determined by the ID in the path and the data used for updating comes from the
// request body.
func UpdateCharacter(c *gin.Context) {
	// TODO: Now name is required
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	existingCharacter := c.MustGet(characterKey).(*characters.Character)

	// The User must be the Character owner in order to update it
	if existingCharacter.UserID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	character := &characters.Character{}
	err := c.Bind(character)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if character.Name == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("name is required"))
		return
	}

	updatedCharacter, err := dbBackend.UpdateCharacter(existingCharacter.ID, character)
	if err != nil {
		log.WithError(err).Error("Failed to update existing character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedCharacter)
}

// DeleteCharacter allows either the Channel owner or the Character owner to delete
// the character.
func DeleteCharacter(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)
	channel := c.MustGet(channelKey).(*channels.Channel)
	character := c.MustGet(characterKey).(*characters.Character)

	// The Channel owner or the User who owns the Character can delete it
	if user.ID != channel.OwnerID && user.ID != character.UserID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	err := dbBackend.DeleteCharacter(character.ID)
	if err != nil {
		log.WithError(err).WithField("characterID", character.ID).Error("Failed to delete character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// LoadCharacter attempts to lookup the Character using the Character ID in the path
// and stores it in the context so the later middleware doesn't have to do it.
func LoadCharacter(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	characterID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		log.WithError(err).Error("Failed to get character id from path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	character, err := dbBackend.GetCharacter(characterID)
	if err != nil {
		if err == characters.ErrCharacterNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		log.WithError(err).WithField("characterID", characterID).Error("Failed look up character.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(characterKey, character)
}
