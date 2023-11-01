package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	address     = "localhost"
	port        = 8080
	loggerLevel = "debug"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    uint64 `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User    string `yaml:"user"`
		DbName  string `yaml:"dbname"`
		Host    string `yaml:"host"`
		Port    uint64 `yaml:"port"`
		SslMode string `yaml:"sslmode"`
		InitDB  struct {
			Init      bool   `yaml:"init"`
			PathToDir string `yaml:"path_to_dir"`
		} `yaml:"init_db"`
	} `yaml:"database"`
	LoggerLvl    string `yaml:"logger_level"`
	RecognizeAPI struct {
		MaxImages    int    `yaml:"max_images"`
		BaseURL      string `yaml:"base_url"`
		CountResults int    `yaml:"count_results"`
		ImageField   string `yaml:"image_field"`
		Token        string `yaml:"token"`
	} `yaml:"recognize_api"`
	TrefleAPI struct {
		BaseURL     string `yaml:"base_url"`
		CountPlants int    `yaml:"count_plants"`
		Token       string `yaml:"token"`
	} `yaml:"trefle_api"`
	PerenualAPI struct {
		BaseURL string `yaml:"base_url"`
		Token   string `yaml:"token"`
	} `yaml:"perenual_api"`
	CookieSettings CookieSettings
}

type CookieSettings struct {
	Secure     bool `yaml:"secure"`
	HttpOnly   bool `yaml:"http_only"`
	ExpireDate struct {
		Years  uint64 `yaml:"years"`
		Months uint64 `yaml:"months"`
		Days   uint64 `yaml:"days"`
	} `yaml:"expire_date"`
}

func New(isInitDB bool) *Config {
	return &Config{
		Server: struct {
			Address string `yaml:"address"`
			Port    uint64 `yaml:"port"`
		}(struct {
			Address string
			Port    uint64
		}{
			Address: address,
			Port:    port,
		}),
		Database: struct {
			User    string `yaml:"user"`
			DbName  string `yaml:"dbname"`
			Host    string `yaml:"host"`
			Port    uint64 `yaml:"port"`
			SslMode string `yaml:"sslmode"`
			InitDB  struct {
				Init      bool   `yaml:"init"`
				PathToDir string `yaml:"path_to_dir"`
			} `yaml:"init_db"`
		}(struct {
			User    string
			DbName  string
			Host    string
			Port    uint64
			SslMode string
			InitDB  struct {
				Init      bool
				PathToDir string
			}
		}{
			User:    "postgres",
			DbName:  "cyber_garden",
			Host:    "localhost",
			Port:    5432,
			SslMode: "disable",
			InitDB: struct {
				Init      bool
				PathToDir string
			}{
				Init:      isInitDB,
				PathToDir: getPwd() + "migrations/plant-database/json",
			},
		}),
		LoggerLvl: loggerLevel,
		RecognizeAPI: struct {
			MaxImages    int    `yaml:"max_images"`
			BaseURL      string `yaml:"base_url"`
			CountResults int    `yaml:"count_results"`
			ImageField   string `yaml:"image_field"`
			Token        string `yaml:"token"`
		}(struct {
			MaxImages    int
			BaseURL      string
			CountResults int
			ImageField   string
			Token        string
		}{
			MaxImages:    5,
			BaseURL:      "https://my-api.plantnet.org/v2/identify/",
			CountResults: 4,
			ImageField:   "images[]",
			Token:        os.Getenv("RECOGNIZE_API_TOKEN"),
		}),
		TrefleAPI: struct {
			BaseURL     string `yaml:"base_url"`
			CountPlants int    `yaml:"count_plants"`
			Token       string `yaml:"token"`
		}(struct {
			BaseURL     string
			CountPlants int
			Token       string
		}{
			BaseURL:     "https://{defaultHost}/api/v1/plants/",
			CountPlants: 5,
			Token:       os.Getenv("TREFLE_API_TOKEN"),
		}),
		PerenualAPI: struct {
			BaseURL string `yaml:"base_url"`
			Token   string `yaml:"token"`
		}(struct {
			BaseURL string
			Token   string
		}{
			BaseURL: "https://{defaultHost}/api/v1/plants/",
			Token:   os.Getenv("PERENUAL_API_TOKEN"),
		}),
		CookieSettings: struct {
			Secure     bool `yaml:"secure"`
			HttpOnly   bool `yaml:"http_only"`
			ExpireDate struct {
				Years  uint64 `yaml:"years"`
				Months uint64 `yaml:"months"`
				Days   uint64 `yaml:"days"`
			} `yaml:"expire_date"`
		}(struct {
			Secure     bool
			HttpOnly   bool
			ExpireDate struct {
				Years  uint64
				Months uint64
				Days   uint64
			}
		}{
			Secure:   true,
			HttpOnly: true,
			ExpireDate: struct {
				Years  uint64
				Months uint64
				Days   uint64
			}{
				Years:  0,
				Months: 0,
				Days:   7,
			},
		}),
	}
}

func (c *Config) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) FormatDbAddr() string {
	return fmt.Sprintf(
		"host=%s user=%s password=admin dbname=%s port=%d sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.DbName,
		c.Database.Port,
		c.Database.SslMode,
	)
}

func ParseFlag(path *string, isInit *bool) {
	flag.StringVar(path, "ConfigPath", "configs/app/local.yaml", "Path to Config")
	flag.BoolVar(isInit, "InitDB", false, "Needs to init DB with plants content")
}

func getPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
