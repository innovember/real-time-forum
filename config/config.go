package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Database struct {
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Path       string `json:"path,omitempty"`
	Driver     string `json:"driver,omitempty"`
	Filename   string `json:"fileName,omitempty"`
	SchemesDir string `json:"schemesDir,omitempty"`
}

type Server struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type Config struct {
	Database Database `json:"database"`
	Server   Database `json:"server"`
}

func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return config, nil
}

func getDBConnString(database Database) string {
	return filepath.Join(database.Path, database.Filename)
}

func (c *Config) GetProdDBConnString() string {
	return getDBConnString(c.Database)
}

func (c *Config) GetDBDriver() string {
	return c.Database.Driver
}

func (c *Config) GetDBSchemesDir() string {
	return c.Database.SchemesDir
}

func (c *Config) GetDBFilename() string {
	return c.Database.Filename
}

func (c *Config) GetDBPath() string {
	return c.Database.Path
}
