package configs

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

type ServerConfiguration struct {
	Port string
}

type DatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}
