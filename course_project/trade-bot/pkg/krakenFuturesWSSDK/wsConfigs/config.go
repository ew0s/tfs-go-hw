package wsConfigs

type WSAPIConfig struct {
	Requests WSAPIRequestsConfig
	Kraken   KrakenWSAPIConfiguration
}

type KrakenWSAPIConfiguration struct {
	WSAPIURL   string
	PublicKey  string
	PrivateKey string
}

type WSAPIRequestsConfig struct {
	WriteWait      int
	PongWait       int
	PingPeriod     int
	MaxMessageSize int
}
