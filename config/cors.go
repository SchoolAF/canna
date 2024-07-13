package config

var allowedOrigins = []string{
	"https://scipio.hlcyn.co/",
	"https://thescipio.github.io/",
}

// GetAllowedOrigins returns the list of allowed origins
func GetAllowedOrigins() []string {
	return allowedOrigins
}
