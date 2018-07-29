package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/andrew-boutin/dndtextapi/messages"

	"github.com/andrew-boutin/dndtextapi/users"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RegisterAdminRoutes adds the admin routes.
func RegisterAdminRoutes(g *gin.RouterGroup) {
	// Routes to admin channels
	g.GET("/test", RequiredHeadersMiddleware(acceptHeader), AdminGetChannels)
	g.GET("/admin/channels/:id", RequiredHeadersMiddleware(acceptHeader), AdminGetChannel)
	g.PUT("/admin/channels/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), AdminUpdateChannel)
	g.DELETE("/admin/channels/:id", AdminDeleteChannel)

	// Routes to admin messages
	g.GET("/admin/messages", RequiredHeadersMiddleware(acceptHeader), AdminGetMessages)
	g.GET("/admin/messages/:id", RequiredHeadersMiddleware(acceptHeader), AdminGetMessage)
	g.PUT("/admin/messages/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), AdminUpdateMessage)
	g.DELETE("/admin/messages/:id", AdminDeleteMessage)

	// Routes to admin users
	g.GET("/admin/users", RequiredHeadersMiddleware(acceptHeader), AdminGetUsers)
	g.GET("/admin/users/:id", RequiredHeadersMiddleware(acceptHeader), AdminGetUser)
	g.PUT("/admin/users/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), AdminUpdateUser)
	g.DELETE("/admin/users/:id", AdminDeleteUser)
}

// RequireAdminHandler requires that the authenticated User be an admin or else
// access is denied.
func RequireAdminHandler(c *gin.Context) {
	user := GetAuthenticatedUser(c)
	if !user.IsAdmin {
		log.Error("Non admin user attempted to access route that requires admin.")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// AdminGetChannels retrieves all of the Channels.
func AdminGetChannels(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	allChannels, err := dbBackend.GetAllChannels(nil)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve all channels.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allChannels)
}

// AdminGetChannel retrieves the Channel matching the id
// in the path.
func AdminGetChannel(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	channel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to retrieve channel.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, channel)
}

// AdminUpdateChannel updates the Channel matching the id
// in the path using the data from the request body.
func AdminUpdateChannel(c *gin.Context) {
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

	updatedChannel, err := dbBackend.UpdateChannel(channelID, channel)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedChannel)
}

// AdminDeleteChannel deletes the Channel matching the id
// in the path.
func AdminDeleteChannel(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dbBackend.DeleteChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to delete channel.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// AdminGetMessages retrieves all of the Messages
// for the Channel matching the required query parameter
// channelID.
func AdminGetMessages(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := QueryParamAsIntExtractor(c, channelIDQueryParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	allMessages, err := dbBackend.GetMessagesInChannel(channelID, nil)
	if err != nil {
		// TODO: Channel not found 404?
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allMessages)
}

// AdminGetMessage retrieves the Message matching the id
// in the path.
func AdminGetMessage(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	message, err := dbBackend.GetMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to retrieve message.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, message)
}

// AdminUpdateMessage updates the Message matching the id
// in the path using the data from the request body.
func AdminUpdateMessage(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedMessage)
}

// AdminDeleteMessage deletes the Message matching the id
// in the path.
func AdminDeleteMessage(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	messageID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dbBackend.DeleteMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to retrieve message.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// AdminGetUsers retrieves all of the Users.
func AdminGetUsers(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	allUsers, err := dbBackend.GetAllUsers()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve all users.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allUsers)
}

// AdminGetUser retrieves the User matching the id in the
// path.
func AdminGetUser(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	userID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := dbBackend.GetUserByID(userID)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

// AdminUpdateUser updates the User matching the id in the path
// with the data from the request body.
func AdminUpdateUser(c *gin.Context) {
	// TODO: Make IsAdmin immutable everywhere in validation (db call doesn't update this so ok for now)
	dbBackend := GetDBBackend(c)

	userID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := &users.User{}
	err = c.Bind(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Look up the existing User so we can see if they're an admin or not
	existingUser, err := dbBackend.GetUserByID(userID)
	if err != nil {
		if err == users.ErrUserNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to look up User before update.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Prevent updates to other admins
	if existingUser.IsAdmin {
		log.WithError(err).Error("Admin attempted to update another admin user.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	updatedUser, err := dbBackend.UpdateUser(userID, user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// AdminDeleteUser deletes the User matching the id in the path.
func AdminDeleteUser(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	userID, err := PathParamAsIntExtractor(c, idPathParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Look up the existing User so we can make sure they're not an admin
	existingUser, err := dbBackend.GetUserByID(userID)
	if err != nil {
		if err == users.ErrUserNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Failed to retrieve user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Prevent deletion of another admin
	if existingUser.IsAdmin {
		log.WithError(err).Error("Admin attempted to delete another admin.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = dbBackend.DeleteUser(userID)
	if err != nil {
		if err == users.ErrUserNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		log.WithError(err).Error("Failed to retrieve user.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
