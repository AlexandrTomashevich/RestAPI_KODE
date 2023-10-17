package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type Config struct {
	ServerPort int      `yaml:"port"`
	Database   Database `yaml:"db"`
}

func (db *Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.Username, db.Password, db.Dbname, db.Sslmode,
	)
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	fmt.Println(config)

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return config, fmt.Errorf("no DB_PASSWORD set in environment variables")
	}
	config.Database.Password = password

	return config, nil
}
