// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/users"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterUsersRoutes registers all of the Users routes with their
// associated middleware.
func RegisterUsersRoutes(g *gin.RouterGroup) {
	g.GET("/users/:id", ValidateHeaders(acceptHeader), GetUser)
	g.PUT("/users/:id", ValidateHeaders(acceptHeader, contentTypeHeader), UpdateUser)
	g.DELETE("/users/:id", DeleteUser)
}

// GetUser retrieves the User matching the id in the path.
func GetUser(c *gin.Context) {
	user := GetAuthenticatedUser(c)

	userIDFromPath, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		log.WithError(err).Error("Failed to get User id from path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Users are only allowed to retrieve themselves
	if userIDFromPath != user.ID {
		log.Error("User attempted to get a User other than themselves.")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser allows a User to update some of their own User data.
func UpdateUser(c *gin.Context) {
	// TODO: Validate only some fields attempted to be updated like bio
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	userIDFromPath, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		log.WithError(err).Error("Failed to get user id from path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Users are only allowed to update themselves
	if userIDFromPath != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Get the provided User data out of the request body
	userBody := &users.User{}
	err = c.Bind(userBody)
	if err != nil {
		log.WithError(err).Error("Issue reading user from request body.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedUser, err := dbBackend.UpdateUser(user.ID, userBody)
	if err != nil {
		log.WithError(err).Error("Failed to update user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser allows a User to delete their own User data.
func DeleteUser(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	userIDFromPath, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Users are only allowed to delete themselves
	if userIDFromPath != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	err = dbBackend.DeleteMessagesFromUser(userIDFromPath)
	if err != nil {
		log.WithError(err).Error("Failed to delete messages from user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = dbBackend.DeleteCharactersFromUser(userIDFromPath)
	if err != nil {
		log.WithError(err).Error("Failed to delete characters from user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = dbBackend.DeleteUser(user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
