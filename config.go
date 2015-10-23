package ice

import (
	_ "log"
	"os"
	"path"
)

type config struct {
	Name, Version             string
	ServerAddress             string
	DBType, DBURL             string
	EmailFrom                 string
	MailgunKey, MailgunDomain string
	MigrationPath             string
}

var Config config

func LoadConfig() {
	var cp = os.Getenv("ICECONFIGPATH")
	if cp == "" {
		cp = ".env"
	} else {
		cp = path.Join(path.Clean(cp), ".env")
	}

	fs, err := os.Open(cp)
	if err != nil {
		panic("Error opening config file " + cp)
	}
	defer fs.Close()

	err = ParseJSON(fs, &Config)
	if err != nil {
		panic("Error parsing config file " + err.Error())
	}
}
