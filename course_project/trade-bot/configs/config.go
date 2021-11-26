package configs

type Configuration struct {
	Server          ServerConfiguration
	Client          ClientConfiguration
	PostgreDatabase PostgreDatabaseConfiguration
	RedisDatabase   RedisDatabaseConfiguration
	Kraken          KrakenConfiguration
}

type ServerConfiguration struct {
	Port string
}

type ClientConfiguration struct {
	URL string
}

type PostgreDatabaseConfiguration struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisDatabaseConfiguration struct {
	Port string
}

type KrakenConfiguration struct {
	APIURL string
}
