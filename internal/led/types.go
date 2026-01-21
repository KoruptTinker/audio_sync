package led

import (
	"github.com/KoruptTinker/audio-sync/internal/config"
)

type BLEWriter interface {
	WriteWithoutResponse(data []byte) error
}

type Controller struct {
	cfg         config.LEDConfig
	conn        BLEWriter
	dataChannel <-chan float64
}

type Manager struct {
	cfg         config.LEDConfig
	bleCfg      config.BluetoothConfig
	controllers []*Controller
}
