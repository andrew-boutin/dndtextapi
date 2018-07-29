// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/users"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterUsersRoutes registers all of the Users routes with their
// associated middleware.
func RegisterUsersRoutes(g *gin.RouterGroup) {
	g.GET("/users", RequiredHeadersMiddleware(acceptHeader), GetUsersInChannel)
	g.GET("/users/:id", RequiredHeadersMiddleware(acceptHeader), GetUser)
	g.PUT("/users/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateUser)
	g.DELETE("/users/:id", DeleteUser)
}

// GetUsersInChannel retrieves all of the Users who are members of the Channel
// matching the required query parameter channelID.
func GetUsersInChannel(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	dbBackend := GetDBBackend(c)

	channelID, err := QueryParamAsIntExtractor(c, "channelID")
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

	// Private Channels require that the User be a member to get the
	// membership list.
	var isMember bool
	if channel.IsPrivate {
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

	usersInChannel, err := dbBackend.GetUsersInChannel(channelID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, usersInChannel)
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
		c.AbortWithStatus(http.StatusUnauthorized)
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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Users are only allowed to update themselves
	if userIDFromPath != user.ID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get the provided User data out of the request body
	userBody := &users.User{}
	err = c.Bind(userBody)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedUser, err := dbBackend.UpdateUser(user.ID, userBody)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// TODO: What about Messages & Channel memberships?

	err = dbBackend.DeleteUser(user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
