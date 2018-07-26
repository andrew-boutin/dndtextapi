// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package configs

// AuthenticationConfiguration holds the authentication configuration data
// that matches the config file.
type AuthenticationConfiguration struct {
	// Accounts is the URL to use to look up User profile data with the authentication client
	Accounts string

	// Oauth2 is the oauth2 URL to use for authentication with the client
	Oauth2 string

	// ID is the client id
	ID string

	// Secret is the client secret
	Secret string
}
