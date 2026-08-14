package config

import (
	"log"
	"sync/atomic"
)

type Store struct {
	value atomic.Value
}

func (s *Store) Load() Config {
	return s.value.Load().(Config)
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
