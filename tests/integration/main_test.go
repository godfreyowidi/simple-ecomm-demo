package integration

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load() // load .env silently
	os.Exit(m.Run())
}
