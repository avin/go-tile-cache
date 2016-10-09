package models

import (
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"time"
)

type Config struct {
	Servers     []ServersConfigItem
	Ttl         int
	TtlDuration time.Duration
	Tries       int
}

type ServersConfigItem struct {
	Alias      string
	Url        string
	Subdomains []string
}

// define the type as a generic map
var config Config
var configLoaded bool

func GetConfig() (*Config) {

	//if config already loaded - use exist copy
	if (configLoaded) {
		return &config
	}

	//read config file
	serversConfigFile := filepath.Join("conf", "config.json")

	file, err := ioutil.ReadFile(serversConfigFile)
	if err != nil {
		fmt.Println("Cannot open servers configuration file:", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &config)

	//convert ttl to ttlDuration
	config.TtlDuration = time.Duration(config.Ttl) * time.Second

	configLoaded = true

	return &config
}