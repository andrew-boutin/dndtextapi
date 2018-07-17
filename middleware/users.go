package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/gin-gonic/gin"
)

// RegisterUsersMiddleware registers all of the Users routes with their
// associated middleware.
func RegisterUsersMiddleware(r *gin.Engine) {
	r.GET("/users", RequiredHeadersMiddleware(acceptHeader), GetUsersInChannel)
}

// GetUsersInChannel retrieves all of the Users who are members of the Channel
// matching the query parameter channelID.
func GetUsersInChannel(c *gin.Context) {
	userID := GetAuthenticatedUserID()
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
	if channel.IsPrivate {
		isMember, err := dbBackend.IsUserInChannel(userID, channelID)
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

	users, err := dbBackend.GetUsersInChannel(channelID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}
