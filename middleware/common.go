package middleware

import (
	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/gin-gonic/gin"
)

const dbBackendKey = "dbBackendKey"

func RegisterMiddleware(r *gin.Engine, backend backends.Backend) {
	r.Use(ContextInjectionMiddleware(backend))

	RegisterChannelsMiddleware(r)
	// RegisterUsersMiddleware(r)
	// RegisterMessagesMiddleware(r)

	/*
		Channels:
		- Create: Authn users can create channels of type "group". A corresponding
		  channel of type "output" is automatically created.
		- Update: Authn users can update channels they're owners of.
		- Delete: Authn users can delete channels they're owners of.
		Messages: List, Get, Create, Update, Delete
		Users: List, Get, Create, Update, Delete
	*/
}

func GetDBBackend(c *gin.Context) backends.Backend {
	return c.MustGet(dbBackendKey).(backends.Backend)
}

func ContextInjectionMiddleware(backend backends.Backend) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbBackendKey, backend)
	}
}
