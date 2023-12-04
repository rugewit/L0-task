package additional

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

func LoadViper(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Could not read config", err)
		return err
	}
	return nil
}

func GetIntVariableFromViper(name string) (int, error) {
	err := LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load viper")
		return -1, err
	}

	variableStr := viper.Get(name).(string)

	variable, err := strconv.Atoi(variableStr)
	if err != nil {
		log.Fatalln(err)
		return -1, err
	}
	return variable, nil
}
