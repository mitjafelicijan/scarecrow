package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var config Config

// Service ...
type Service struct {
	Path  string `yaml:"path"`
	Proxy string `yaml:"proxy"`
}

// Config ...
type Config struct {
	Listen          string    `yaml:"listen"`
	Metrics         bool      `yaml:"metrics"`
	GZIP            bool      `yaml:"gzip"`
	Logging         bool      `yaml:"logging"`
	Verbose         bool      `yaml:"verbose"`
	Heartbeat       bool      `yaml:"heartbeat"`
	Console         bool      `yaml:"console"`
	TLS             bool      `yaml:"tls"`
	Domain          string    `yaml:"domain"`
	Email           string    `yaml:"email"`
	ServiceRegistry []Service `yaml:"registry"`
}

func parseConfigFile(filename string) Config {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Error reading config file")
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println("Error parsing config file")
		os.Exit(1)
	}

	return config
}
