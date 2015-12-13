package ice

import (
	"os"
path 	"path/filepath"
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
	Production                bool
	Settings                  map[string]string
}

func (c *config) Get(key string) string {
	return c.Settings[key]
}

var Config config

func findEnv(dir string)string{
fn:=path.Join(dir,".env")
_,err:=os.Lstat(fn)
if err!=nil{
newDir:=path.Dir(dir)
if dir == newDir{
panic("Not inside an ice directory")
}
return findEnv(newDir)
}
return fn
}

func LoadConfig() {
	var cp = os.Getenv("ICECONFIGPATH")
	if cp == "" {
dir,err:=os.Getwd()
if err!=nil{
panic(err)
}
cp=findEnv(dir)
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

	if Config.Settings == nil {
		Config.Settings = make(map[string]string)
	}
}
