package server

type ServerConfig struct {
	Host string
	Port string
	RootDir string
	HttpsEnabled bool
	CertPath string
	KeyPath string
	CorsEnabled bool
}
