package configHelper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	AppSettings       map[string]interface{} `json:"appSettings"`
	ConnectionStrings map[string]interface{} `json:"connectionStrings"`
}

func GetConfigs(jsonFilePath string) *Configuration {
	f, _ := os.Open(jsonFilePath)
	defer f.Close()

	configuration := Configuration{}
	b, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(b, &configuration)
	if err != nil {
		log.Fatal(err)
	}
	return &configuration
}
