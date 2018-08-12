// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/bots"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterBotsRoutes registers all of the bot routes with their
// associated middleware
func RegisterBotsRoutes(g *gin.RouterGroup) {
	g.GET("/bots", ValidateHeaders(acceptHeader), GetBots)
	g.POST("/bots", ValidateHeaders(acceptHeader, contentTypeHeader), CreateBot)
	g.GET("/bots/:botID", ValidateHeaders(acceptHeader), GetBot)
	g.PUT("/bots/:botID", ValidateHeaders(acceptHeader, contentTypeHeader), UpdateBot)
	g.DELETE("/bots/:botID", DeleteBot)

	g.GET("/bots/:botID/creds", ValidateHeaders(acceptHeader), GetBotCreds)
}

// GetBots retrieves all bots.
func GetBots(c *gin.Context) {
	db := GetDBBackend(c)

	bs, err := db.GetAllBots()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve all bots.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bs)
}

// GetBot retrieves the bot that matches the id from the path.
func GetBot(c *gin.Context) {
	db := GetDBBackend(c)

	id, err := PathParamAsIntExtractor(c, "botID")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	b, err := db.GetBot(id)
	if err != nil {
		if err == bots.ErrBotNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).WithField("id", id).Error("Failed to retrieve bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, b)
}

// CreateBot creates a new bot with the given data from the request body.
func CreateBot(c *gin.Context) {
	// TODO: Validation on input
	db := GetDBBackend(c)

	botInput := &bots.Bot{}
	err := c.Bind(botInput)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	b, err := db.CreateBot(botInput)
	if err != nil {
		log.WithError(err).Error("Failed to create bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create client credentials for the bot
	_, err = db.CreateBotCreds(b)
	if err != nil {
		log.WithError(err).WithField("id", b.ID).Error("Failed to create bot credentials.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, b)
}

// UpdateBot allows the user, if they're the bot owner, to update the
// bot matching the id in the path with the data in the request body.
func UpdateBot(c *gin.Context) {
	// TODO: Can only update workspace
	db := GetDBBackend(c)
	user := GetAuthenticatedUser(c)

	botInput := &bots.Bot{}
	err := c.Bind(botInput)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := PathParamAsIntExtractor(c, "botID")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingBot, err := db.GetBot(id)
	if err != nil {
		if err == bots.ErrBotNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).WithField("id", id).Error("Failed to look up existing bog.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Only the bot owner can update the bot
	if existingBot.OwnerID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	b, err := db.UpdateBot(id, botInput)
	if err != nil {
		log.WithError(err).Error("Failed to update bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, b)
}

// DeleteBot allows the user, if they're the bot owner, to delete the
// bot matching the id in the path.
func DeleteBot(c *gin.Context) {
	db := GetDBBackend(c)
	user := GetAuthenticatedUser(c)

	id, err := PathParamAsIntExtractor(c, "botID")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingBot, err := db.GetBot(id)
	if err != nil {
		if err == bots.ErrBotNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).WithField("id", id).Error("Failed to look up existing bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Only the bot owner can delete the bot
	if existingBot.OwnerID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	err = db.DeleteBotCreds(id)
	if err != nil {
		log.WithError(err).WithField("botID", id).Error("Failed to delete bot credentials.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = db.DeleteBot(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Failed to delete existing bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBotCreds allows the user, if they're the bot owner, to get the
// credentials for the bot matching the id in the path.
func GetBotCreds(c *gin.Context) {
	db := GetDBBackend(c)
	user := GetAuthenticatedUser(c)

	id, err := PathParamAsIntExtractor(c, "botID")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existingBot, err := db.GetBot(id)
	if err != nil {
		if err == bots.ErrBotCredsNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).WithField("id", id).Error("Failed to look up existing bot.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Only the bot owner can get the credentials for the bot
	if existingBot.OwnerID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	creds, err := db.GetBotCreds(id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Failed to get bot credentials.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, creds)
}
