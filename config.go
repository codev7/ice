package ice

import (
	_ "log"
	"os"
)

type config struct {
	ServerAddress                      string
	DBHost, DBUser, DBPassword, DBName string
	EmailFrom                          string
	MailgunKey, MailgunDomain          string
}

var Config config

func loadConfig() {
	fs, err := os.Open(".env")
	if err != nil {
		fs, err = os.Open("../../.env")
		if err != nil {
			panic("Error opening config file " + err.Error())
		}
	}
	defer fs.Close()

	err = ParseJSON(fs, &Config)
	if err != nil {
		panic("Error parsing config file " + err.Error())
	}
}
