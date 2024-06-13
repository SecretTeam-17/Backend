package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Структура конфига
type Config struct {
	Env           string `yaml:"env" env-default:"prod"`
	StoragePath   string `yaml:"storage_path" env-required:"true"`
	StoragePasswd string `yaml:"storage_passwd" env-default:"DB_PASSWD"`
	HTTPServer    `yaml:"http_server"`
}
type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"localhost:80"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad - инициализируем данные из конфига. Если не смогли, то завершаем приложение с ошибкой.
func MustLoad() *Config {
	// Берем путь к файлу конфига из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")

	// Вариант для локала
	// configPath := "./configs/local.yaml"

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Проверяем, существует ли файл конфига
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// Читаем из файла конфига в структуру и возвращаем ее
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
