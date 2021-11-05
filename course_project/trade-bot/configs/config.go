package configs

type Configuration struct {
	Server ServerConfiguration
}

type ServerConfiguration struct {
	Port string
}
