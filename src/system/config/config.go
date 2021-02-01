package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ferrany1/log2file/src/directory"
	"io/ioutil"
)

var (
	CF  *Config
	err error
)

type Config struct {
	DB     postgresDB `json:"postgresDB"`
	Router ginRouter  `json:"ginRouter"`
}

type postgresDB struct {
	Domain   string         `json:"domain"`
	Username string         `json:"user"`
	Password string         `json:"password"`
	DBName   string         `json:"dbname"`
	Port     int            `json:"port"`
	Tables   postgresTables `json:"tables"`
}

type postgresTables struct {
	Users    string `json:"users"`
	Balances string `json:"balances"`
}

type ginRouter struct {
	Port int `json:"port"`
}

// Reads config file
func ReadConfig(filename string) error {
	CF, err = readConfigFile(filename)
	return err
}

// Reads file from directory
func readConfigFile(filename string) (cf *Config, err error) {
	fi, dir, err := directory.ReadCurrentDirectory()
	if err != nil {
		return cf, err
	}

	for _, fn := range fi {
		if fn.Name() == filename {
			b, err := ioutil.ReadFile(dir + "/" + filename)
			if err != nil {
				return cf, errors.New(fmt.Sprintf("failed read file. %s", err.Error()))
			}

			err = json.Unmarshal(b, &cf)
			if err != nil {
				return cf, errors.New(fmt.Sprintf("failed to unmarshal file to type. %s", err.Error()))
			}
			return cf, nil
		}
	}
	return cf, errors.New(fmt.Sprintf("no config found in dir."))
}