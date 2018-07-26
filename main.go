// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package main

import (
	"github.com/andrew-boutin/dndtextapi/backends"
	"github.com/andrew-boutin/dndtextapi/configs"
	"github.com/andrew-boutin/dndtextapi/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Read in config
	configuration := configs.LoadConfig()

	// Initialize backend
	backend := backends.InitBackend(configuration.Backend)

	// Initalize authentication data
	middleware.InitAuthentication(configuration.Authentication)

	// Set up server
	r := gin.Default()
	middleware.RegisterMiddleware(r, backend)
	r.Run(":8080")
}
