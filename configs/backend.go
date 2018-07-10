package configs

// BackendConfiguration holds the backend configuration data
// that matches the config file.
type BackendConfiguration struct {
	Type   string
	User   string
	PW     string
	DBName string
}
