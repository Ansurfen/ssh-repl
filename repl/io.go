package repl

import (
	"os"

	"github.com/spf13/viper"
)

func GetConf(confName, confType, dir string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigName(confName)
	conf.SetConfigType(confType)
	conf.AddConfigPath(dir)
	Panic(conf.ReadInConfig())
	return conf
}

func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_APPEND, 0666)
}
