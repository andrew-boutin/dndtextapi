package middleware

import (
	"fmt"
	"net/http"
	"strconv"

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

// Errors
var (
	// ErrParamNotInPath is the error to use when a parameter is expected but
	// not found in the path.
	ErrParamNotInPath = fmt.Errorf("expected path parameter not found")

	// ErrParamNotInt is the error to use when a parameter in the path is
	// expected to be an integer but isn't.
	ErrParamNotInt = fmt.Errorf("expected path parameter to be an integer")
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
				c.AbortWithError(http.StatusBadRequest, c.Error(fmt.Errorf("missing required header %s", requiredHeader)))
				return
			}

			switch requiredHeader {
			case acceptHeader:
				if val != applicationJSONHeaderVal {
					// TODO: 415 Unsupported Media Type?
					log.WithField(acceptHeader, val).Error("Invalid header value.")
					c.AbortWithError(http.StatusBadRequest, c.Error(fmt.Errorf("invalid %s header value %s", acceptHeader, val)))
					return
				}
			case contentTypeHeader:
				if val != applicationJSONHeaderVal {
					// TODO: 415 Unsupported Media Type?
					log.WithField(contentTypeHeader, val).Error("Invalid header value.")
					c.AbortWithError(http.StatusBadRequest, c.Error(fmt.Errorf("invalid %s header value %s", contentTypeHeader, val)))
					return
				}
			}
		}
	}
}

// PathParamExtractor extracts a parameter from the gin.Context by using
// the given name. If no parameter is found then an error is returned.
func PathParamExtractor(c *gin.Context, name string) (string, error) {
	p := c.Param(name)

	if len(p) <= 0 {
		return "", ErrParamNotInPath
	}

	return p, nil
}

// PathParamAsIntExtractor extracts a parameter from the gin.Context
// by using the given name and returns it as an integer.
func PathParamAsIntExtractor(c *gin.Context, name string) (int, error) {
	pStr, err := PathParamExtractor(c, name)

	if err != nil {
		return 0, err
	}

	pInt, err := strconv.Atoi(pStr)

	if err != nil {
		return 0, ErrParamNotInt
	}

	return pInt, nil
}
