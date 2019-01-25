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
	Path  string `yaml:"path" json:"path"`
	Proxy string `yaml:"proxy" json:"proxy"`
}

// Timeout ...
type Timeout struct {
	Read  int `yaml:"read" json:"read"`
	Write int `yaml:"write" json:"write"`
}

// Config ...
type Config struct {
	Listen          string    `yaml:"listen" json:"listen"`
	Metrics         bool      `yaml:"metrics" json:"metrics"`
	GZIP            bool      `yaml:"gzip" json:"gzip"`
	Logging         bool      `yaml:"logging" json:"logging"`
	Verbose         bool      `yaml:"verbose" json:"verbose"`
	Heartbeat       bool      `yaml:"heartbeat" json:"heartbeat"`
	Console         bool      `yaml:"console" json:"console"`
	TLS             bool      `yaml:"tls" json:"tls"`
	Domain          string    `yaml:"domain" json:"domain"`
	Email           string    `yaml:"email" json:"email"`
	StatsDBFreq     int       `yaml:"stats-db-freq" json:"statsDbFreq"`
	Timeout         Timeout   `yaml:"timeout" json:"timeout"`
	ServiceRegistry []Service `yaml:"registry" json:"registry"`
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
