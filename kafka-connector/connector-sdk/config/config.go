package config

import (
	"fmt"
	"io/ioutil"
	//"os"
	//"strings"

	yaml "gopkg.in/yaml.v2"
)

func NewServices() *Services {
	var services Services
	services.Openfx.FxGatewayURL = DefaultGatewayURL

	services.Functions = make(map[string]Function, 0)

	return &services
}

func ParseConfigFile(file string) (*Services, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var services Services
	err = yaml.Unmarshal(fileData, &services)
	if err != nil {
		fmt.Printf("Error with YAML Config file\n")
		return nil, err
	}

	return &services, nil
}

func GetFxGatewayURL() string {
	var url string

	url = DefaultGatewayURL

	return url
}
