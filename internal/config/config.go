package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func Load() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use defaults
			fmt.Println("Config file not found, using defaults. Create config/config.yaml from config.example.yaml")
		} else {
			return Config{}, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode config: %w", err)
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("audio.devicename", "BlackHole")
	viper.SetDefault("audio.samplerate", 86000)
	viper.SetDefault("audio.channels", 1)

	viper.SetDefault("bluetooth.devicename", "GATT--DEMO")
	viper.SetDefault("bluetooth.devicecount", 2)
	viper.SetDefault("bluetooth.characteristicuuid", "0000fff3-0000-1000-8000-00805f9b34fb")

	viper.SetDefault("led.updateinterval", 35*time.Millisecond)
	viper.SetDefault("led.minbrightness", 100.0)
	viper.SetDefault("led.maxbrightness", 900.0)
	viper.SetDefault("led.threshold", 0.01)
	viper.SetDefault("led.gain", 1.8)
	viper.SetDefault("led.curveexponent", 4.0)
}

func Default() Config {
	setDefaults()
	var cfg Config
	viper.Unmarshal(&cfg)
	return cfg
}
