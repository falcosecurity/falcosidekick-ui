package configuration

type Configuration struct {
	// DisplayMode   string `json:"display-mode"`
	ListenAddress string `json:"listen-address"`
	ListenPort    int    `json:"listen-port"`
	RedisServer   string `json:"redis-server"`
	DevMode       bool   `json:"dev-mode"`
}

var config *Configuration

func CreateConfiguration() *Configuration {
	config = new(Configuration)
	return config
}

func GetConfiguration() *Configuration {
	return config
}