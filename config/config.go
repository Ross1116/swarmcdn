package config

type Config struct {
	ChunkSize int
	ServerURL string
}

var AppConfig *Config

func InitConfig() {
	AppConfig = &Config{
		ChunkSize: 512 * 1024, // 512 KB
		ServerURL: "http://localhost:8080",
	}
}
