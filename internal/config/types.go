package config

import (
	"time"
)

type Config struct {
	Audio     AudioConfig
	Bluetooth BluetoothConfig
	LED       LEDConfig
}

type AudioConfig struct {
	DeviceName string
	SampleRate uint32
	Channels   uint32
}

type BluetoothConfig struct {
	DeviceName         string
	DeviceCount        int
	CharacteristicUUID string
}

type LEDConfig struct {
	UpdateInterval time.Duration
	MinBrightness  float64
	MaxBrightness  float64
	Threshold      float64
	Gain           float64
	CurveExponent  float64
}
