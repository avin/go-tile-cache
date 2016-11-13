package models

import (
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"time"
	"log"
)

type Config struct {
	Servers       []ServersConfigItem
	Ttl           int
	TtlDuration   time.Duration
	Tries         int
	Proxy         string
	Cache         string
	ClearOldCache bool
	Mongodb       MongodbConfig
}

type MongodbConfig struct {
	Host   string
	DbName string
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	serversConfigFile := filepath.Join(dir, "conf", "config.json")

	file, err := ioutil.ReadFile(serversConfigFile)
	if err != nil {
		fmt.Println("Cannot open servers configuration file:", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &config)

	//convert ttl to ttlDuration
	config.TtlDuration = time.Duration(config.Ttl) * time.Second

	//setup proxy
	if (len(config.Proxy) > 0) {
		os.Setenv("HTTP_PROXY", config.Proxy)
	}

	configLoaded = true

	return &config
}