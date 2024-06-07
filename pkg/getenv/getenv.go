package getenv

import (
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

func GetEnvConfig(config string) string {
	envVar := os.Getenv(config)
	if envVar == "" {
		err := gotenv.Load(".env")
		if err != nil {
			panic(fmt.Sprintf("environment variable %s was not found.", config))
		}
		envVar = os.Getenv(config)
	}
	if config == "" {
		panic("environment config not found")
	}
	return envVar
}
