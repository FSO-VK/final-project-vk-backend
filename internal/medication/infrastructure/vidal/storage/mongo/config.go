package mongo

// Config is a configuration for MongoDB storage.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Log      bool
}
