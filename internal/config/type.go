package config

// Config Storing all the configurations inside this variable
type Config struct {
	Token     string `json:"token"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	SSLMode   string `json:"ssl_mode"`
	Database  string `json:"database"`
	Migration bool   `json:"migration"`
}
