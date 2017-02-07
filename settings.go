package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Settings struct {
	LoginToken string `yaml:"login_token"`
	Domain     string `yaml:"domain"`
	SubDomain  string `yaml:"sub_domain"`
	Type  string `yaml:"type"`
	IPService  string `yaml:"ip_service"`
	Ttl        int `yaml:"ttl"`
}

func Load(path string, settings *Settings) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		return err
	}

	err = yaml.Unmarshal(data, settings)
	log.Println(settings)
	if err != nil {
		fmt.Println("Error occurs while unmarshal config file, please make sure config file correct!")
		return err
	}

	return nil
}