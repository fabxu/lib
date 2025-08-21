package config

import (
	"os"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func Load(cfgFilePath string) error {
	// load envs form local file .env/conf
	loadEnvsFromLocal()

	viper.SetConfigFile(cfgFilePath)

	return viper.ReadInConfig()
}

func Global() *viper.Viper {
	return viper.GetViper()
}

func loadEnvsFromLocal() {
	filename := ".env/conf"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	envs, err := gotenv.StrictParse(file)
	if err != nil {
		panic(err)
	}

	for key, value := range envs {
		if _, ok := os.LookupEnv(key); ok {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}
