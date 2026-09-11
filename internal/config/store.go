package config

import (
	"log"
	"os"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

type Store struct {
	value atomic.Value
}

func (s *Store) Load() *Config {
	return s.value.Load().(*Config)
}

func (s *Store) Reload() error {
	cfg, err := LoadConfig("config.yml")
	if err != nil {
		return err
	}

	s.value.Store(cfg)

	log.Println("Config reloaded")

	return nil
}

func (s *Store) UpdateUserInfo(token, name, staffNo, email string) error {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	cfg.Token = token
	cfg.Name = name
	cfg.StaffNo = staffNo
	cfg.Email = email

	data, err = yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile("config.yml", data, 0644); err != nil {
		return err
	}

	return s.Reload()
}

func (s *Store) UpdateUserConfig(token, name, staffNo, email string) error {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	cfg.Token = token
	cfg.Name = name
	cfg.StaffNo = staffNo
	cfg.Email = email

	data, err = yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile("config.yml", data, 0644); err != nil {
		return err
	}

	return s.Reload()
}
