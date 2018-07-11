package middleware

import (
	"net/http"

	"github.com/andrew-boutin/dndtextapi/messages"
	"github.com/gin-gonic/gin"
)

// RegisterMessagesMiddleware registers all of the Message routes with their
// associated middleware
func RegisterMessagesMiddleware(r *gin.Engine) {
	r.GET("/messages", RequiredHeadersMiddleware(acceptHeader), GetMessages)
	r.POST("/messages", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), CreateMessage)
	r.GET("/messages/:id", RequiredHeadersMiddleware(acceptHeader), GetMessage)
	r.PUT("/messages/:id", RequiredHeadersMiddleware(acceptHeader, contentTypeHeader), UpdateMessage)
	r.DELETE("/messages/:id", DeleteMessage)
}

// GetMessages retrieves all of the messages for a given Channel defined by a
// query parameter.
func GetMessages(c *gin.Context) {
	// TODO: If channel is private then authn user must be a member of the channel
	dbBackend := GetDBBackend(c)

	channelIDStr, err := QueryParamAsIntExtractor(c, "channelID")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	messages, err := dbBackend.GetMessagesInChannel(channelIDStr)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, messages)
}

// GetMessage retrieves a single Message using the Message ID
// in the path.
func GetMessage(c *gin.Context) {
	// TODO: User must have access to the channel
	dbBackend := GetDBBackend(c)
	messageID, err := PathParamAsIntExtractor(c, "id")
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

	c.JSON(http.StatusOK, message)
}

// CreateMessage creates a new Message using the data in the
// request body.
func CreateMessage(c *gin.Context) {
	// TODO: Set the userID based on the authn User
	// TODO: User must have access to the channel
	// TODO: Validation - content not empty, etc.
	dbBackend := GetDBBackend(c)
	message := &messages.Message{}
	err := c.Bind(message)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createdMessage, err := dbBackend.CreateMessage(message)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdMessage)
}

// DeleteMessage deletes the message matching the ID in the path.
func DeleteMessage(c *gin.Context) {
	// TODO: Authn user must be the user on the Message
	dbBackend := GetDBBackend(c)
	messageID, err := PathParamAsIntExtractor(c, "id")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = dbBackend.DeleteMessage(messageID)
	if err != nil {
		if err == messages.ErrMessageNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateMessage updates the Message using the ID from the path with
// the data from the request body.
func UpdateMessage(c *gin.Context) {
	dbBackend := GetDBBackend(c)
	messageID, err := PathParamAsIntExtractor(c, "id")
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
		if err == messages.ErrMessageNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedMessage)
}
