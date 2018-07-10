package middleware

import (
	"fmt"

	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// Headers
	contentTypeHeader = "content-type"
	acceptHeader      = "accept"

	// Context keys
	dbBackendKey = "dbBackendKey"

	// Other
	applicationJSONHeaderVal = "application/json"
)

// RegisterMiddleware handles registering all common middleware
// and registering all of the various route groups.
func RegisterMiddleware(r *gin.Engine, backend backends.Backend) {
	r.Use(ContextInjectionMiddleware(backend))

	RegisterChannelsMiddleware(r)
	// RegisterUsersMiddleware(r)
	// RegisterMessagesMiddleware(r)

	/*
		Channels:
		- Create: Authn users can create channels of type "group". A corresponding
		  channel of type "output" is automatically created.
		Messages: List, Get, Create, Update, Delete
		Users: List, Get, Create, Update, Delete
	*/
}

// GetDBBackend pulls the db backend out of the context that
// was previously injected.
func GetDBBackend(c *gin.Context) backends.Backend {
	return c.MustGet(dbBackendKey).(backends.Backend)
}

// ContextInjectionMiddleware injects various data into the context
// so that it will be available throughout the rest of the middleware
// that executes on the route.
func ContextInjectionMiddleware(backend backends.Backend) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbBackendKey, backend)
	}
}

// RequiredHeadersMiddleware is a gin.HandlerFunc wrapper that takes in required
// headers. It returns a gin.HandlerFunc that verifies that the request contains
// all of the required headers and that specific headers have the correct values.
func RequiredHeadersMiddleware(expectedHeaders ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, requiredHeader := range expectedHeaders {
			val := c.GetHeader(requiredHeader)
			if len(val) <= 0 {
				log.WithField("header", requiredHeader).Error("Missing required header.")
				c.Error(fmt.Errorf("missing required header %s", requiredHeader))
				return
			}

			switch requiredHeader {
			case acceptHeader:
				if val != applicationJSONHeaderVal {
					log.WithField(acceptHeader, val).Error("Invalid header value.")
					c.Error(fmt.Errorf("invalid %s header value %s", acceptHeader, val))
					return
				}
			case contentTypeHeader:
				if val != applicationJSONHeaderVal {
					log.WithField(contentTypeHeader, val).Error("Invalid header value.")
					c.Error(fmt.Errorf("invalid %s header value %s", contentTypeHeader, val))
					return
				}
			}
		}
	}
}
