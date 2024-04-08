package server

type ServerConfig struct {
	Host         string
	Port         string
	RootDir      string
	CorsEnabled  bool
	DirViewTheme string
	HttpsEnabled bool
	CertPath     string
	KeyPath      string
}
