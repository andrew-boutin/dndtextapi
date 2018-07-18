// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package configs

// BackendConfiguration holds the backend configuration data
// that matches the config file.
type BackendConfiguration struct {
	Type   string
	User   string
	PW     string
	DBName string
}
