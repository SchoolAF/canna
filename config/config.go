package config

import (
	"os"
)

var JWTSecret = []byte(os.Getenv("AUTH_SECRET"))
