package config

import (
	"bufio"
	"os"
)

func GetSecret(env string) string {
	secret := os.Getenv(env)
	if secret == "" {
		file, err := os.Open(os.Getenv(env + "_FILE"))
		if err != nil {
			return secret
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanned := scanner.Scan(); scanned {
			secret = scanner.Text()
		}
	}
	return secret
}
