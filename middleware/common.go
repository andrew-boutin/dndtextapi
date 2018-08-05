// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/andrew-boutin/dndtextapi/channels"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// Headers
	contentTypeHeader = "content-type"
	acceptHeader      = "accept"

	// Context keys
	dbBackendKey = "dbBackendKey"
	channelKey   = "channelKey"
	characterKey = "characterKey"

	// Other
	applicationJSONHeaderVal = "application/json"
	anyMedia                 = "*/*"
	idPathParam              = "id"
	channelIDPathParam       = "channelID"
)

// Query parameters and their valid values
const (
	// TODO: Find a better way to set these up
	// msgTypeQueryParam can be `story` or `meta`.
	msgTypeQueryParam = "msgType"
	metaMsgType       = "meta"
	storyMsgType      = "story"

	// levelQueryParam can be `member` or `owner`.
	levelQueryParam = "level"
	ownerLevel      = "owner"
	memberLevel     = "member"
)

// Errors used throughout the middleware.
var (
	// ErrPathParamNotFound is the error to use when a parameter is expected but
	// not found in the path.
	ErrPathParamNotFound = fmt.Errorf("expected path parameter not found")

	// ErrPathParamNotInt is the error to use when a parameter in the path is
	// expected to be an integer but isn't.
	ErrPathParamNotInt = fmt.Errorf("expected path parameter to be an integer")

	// ErrQueryParamNotFound is the error to use when a parameter is expected
	// but not found in the query string.
	ErrQueryParamNotFound = fmt.Errorf("expected query parameter not found")

	// ErrQueryParamNotInt
	ErrQueryParamNotInt = fmt.Errorf("expected query parameter to be an integer")
)

// RegisterMiddleware handles registering all common middleware
// and registering all of the various route groups.
func RegisterMiddleware(r *gin.Engine, backend backends.Backend) {
	// TODO: Is it possible to register a middleware at the beginning of all PUT/GET etc. routes?
	r.Use(ContextInjectionMiddleware(backend))

	RegisterAnonymousRoutes(r)

	RegisterAuthenticationRoutes(r)

	authorized := r.Group("/")
	authorized.Use(AuthenticationMiddleware)

	RegisterChannelsRoutes(authorized)
	RegisterUsersRoutes(authorized)
	RegisterMessagesRoutes(authorized)
	RegisterCharactersRoutes(authorized)

	// Set up all of the admin only routes
	admin := authorized.Group("/") // TODO: want this to be `/admin`
	admin.Use(RequireAdminHandler)
	RegisterAdminRoutes(admin)
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

// ValidateHeaders is a gin.HandlerFunc wrapper that takes in headers and returns
// a gin.HandlerFunc that verifies if the headers are present in the request that
// they are valid values.
func ValidateHeaders(headers ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, header := range headers {
			val := c.GetHeader(header)
			if len(val) > 0 {
				switch header {
				case acceptHeader:
					if val != applicationJSONHeaderVal && val != anyMedia {
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
}

// PathParamExtractor extracts a path parameter from the gin.Context by using
// the given name. If no parameter is found then an error is returned.
func PathParamExtractor(c *gin.Context, name string) (string, error) {
	p := c.Param(name)

	if len(p) <= 0 {
		return "", ErrPathParamNotFound
	}

	return p, nil
}

// PathParamAsIntExtractor extracts a path parameter from the gin.Context
// by using the given name and returns it as an integer.
func PathParamAsIntExtractor(c *gin.Context, name string) (int, error) {
	pStr, err := PathParamExtractor(c, name)

	if err != nil {
		return 0, err
	}

	pInt, err := strconv.Atoi(pStr)

	if err != nil {
		return 0, ErrPathParamNotInt
	}

	return pInt, nil
}

// QueryParamExtractor extracts a query parameter from the gin.Context by using
// the given name.
func QueryParamExtractor(c *gin.Context, name string) (string, error) {
	p := c.Query(name)

	if p == "" {
		return "", ErrQueryParamNotFound
	}

	return p, nil
}

// QueryParamAsIntExtractor extracts a query parameter from the gin.Context
// by using the given name and returns it as an integer.
func QueryParamAsIntExtractor(c *gin.Context, name string) (int, error) {
	pStr, err := QueryParamExtractor(c, name)

	if err != nil {
		return 0, err
	}

	pInt, err := strconv.Atoi(pStr)

	if err != nil {
		return 0, ErrQueryParamNotInt
	}

	return pInt, nil
}

// LoadChannelFromPathID attempts to lookup the Channel using the Channel ID in the path and stores it
// in the context so the later middleware doesn't have to do it.
func LoadChannelFromPathID(c *gin.Context) {
	dbBackend := GetDBBackend(c)

	channelID, err := PathParamAsIntExtractor(c, channelIDPathParam)
	if err != nil {
		log.WithError(err).Error("Failed to get channel id from path.")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	channel, err := dbBackend.GetChannel(channelID)
	if err != nil {
		if err == channels.ErrChannelNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		log.WithError(err).WithField("channelID", channelID).Error("Failed look up channel.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(channelKey, channel)
}
