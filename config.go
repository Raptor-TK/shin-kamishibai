package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds server config
type Config struct {
	Path         string   `json:"path,omitempty"`     // config file path
	IP           string   `json:"ip"`                 // network ip interface to listen to
	Port         int      `json:"port"`               // server port
	DBPath       string   `json:"db_path"`            // where db file is stored
	Username     string   `json:"username"`           // username for the http authentication
	Password     string   `json:"password,omitempty"` // one time, and it will be cleared after computed
	Iterations   int      `json:"iterations"`         // safety, min 100,000
	Salt         string   `json:"salt"`               // salt for the crypt
	Crypt        string   `json:"crypt"`              // password hash
	AllowedDirs  []string `json:"allowed_dirs"`       // directory allowed to be browse
	ImageResize  bool     `json:"image_resize"`       // resize images in reader
	ImageQuality int      `json:"image_quality"`      // image quality for resized image
}

// ConfigHashIterations how many times the password should be hashed
const ConfigHashIterations = 100000

// Read read and parse configuration file
func (cfg *Config) Read(fpath string) error {
	byteDat, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteDat, &cfg)
	if err != nil {
		return err
	}

	// sanity check
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("invalid port number " + strconv.Itoa(cfg.Port))
	}
	if cfg.Crypt == "" && len(cfg.Password) < 6 {
		return errors.New("password too short, min of 6")
	}
	if len(cfg.Username) < 3 {
		return errors.New("username too short, min length 3")
	}

	// overwrite
	cfg.Path = fpath
	cfg.Iterations = ConfigHashIterations

	if cfg.DBPath == "" {
		cfg.DBPath = filepath.Dir(fpath) + "/db.txt"
	}

	// hash password
	if cfg.Crypt == "" {
		// generate salt, longer because of limited character list
		cfg.Salt = GenerateString(128)
		// calc password hash
		cfg.Crypt = SHA256Iter(cfg.Password, cfg.Salt, ConfigHashIterations)
		// clear password
		cfg.Password = ""
		// save new cfg file
		err := cfg.Save(cfg.Path)
		if err != nil {
			fmt.Println("failed to save config file (b)")
			return err
		}
	}

	// create thumbnail cache dir if not exists
	cacheDir := filepath.Join(filepath.Dir(cfg.Path), "cache")
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Save save config to json file
func (cfg *Config) Save(fpath string) error {
	// create a copy
	config2 := cfg
	// clear, setup setting
	config2.Path = ""
	config2.Password = ""
	config2.Iterations = ConfigHashIterations

	// save to file
	byteDat2, err := json.MarshalIndent(config2, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fpath, byteDat2, 0644)
	if err != nil {
		return err
	}

	return nil
}
