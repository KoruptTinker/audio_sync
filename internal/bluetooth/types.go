package bluetooth

import (
	"github.com/KoruptTinker/audio-sync/internal/config"
	"tinygo.org/x/bluetooth"
)

type Scanner struct {
	cfg     config.BluetoothConfig
	adapter *bluetooth.Adapter
}

type Connection struct {
	device         bluetooth.Device
	characteristic bluetooth.DeviceCharacteristic
	cfg            config.BluetoothConfig
}
