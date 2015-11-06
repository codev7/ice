package ice

import (
	"os"
	"path"
)

type config struct {
	ConfigRoot                string
	Name, Version             string
	ServerAddress             string
	DBType, DBURL             string
	EmailFrom                 string
	MailgunKey, MailgunDomain string
	MigrationPath             string
	Secret                    string
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
		panic("Error opening config file " + cp + err.Error())
	}
	defer fs.Close()

	err = ParseJSON(fs, &Config)
	if err != nil {
		panic("Error parsing config file " + err.Error())
	}
	Config.ConfigRoot = path.Dir(cp)
	if !path.IsAbs(Config.MigrationPath) {
		Config.MigrationPath = path.Clean(path.Join(Config.ConfigRoot, Config.MigrationPath))
	}

}
