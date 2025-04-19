//go:build dev
// +build dev

package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("HTTP_PORT", "8090")
	os.Setenv("GIN_MODE", gin.DebugMode)
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "5434")
	os.Setenv("DB_USER", "user_challenge_svc")
	os.Setenv("DB_PASSWORD", "user_challenge_svc")
	os.Setenv("DB_NAME", "user_challenge_svc")
	os.Setenv("DB_MAX_CONNECTIONS", "100")
	os.Setenv("DB_SSL", "disable")
}
