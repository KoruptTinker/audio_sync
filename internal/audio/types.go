package audio

import (
	"github.com/KoruptTinker/audio-sync/internal/config"
	"github.com/gen2brain/malgo"
)

type RMSHandler struct {
	cfg      config.AudioConfig
	channels []chan float64
	ctx      *malgo.AllocatedContext
	device   *malgo.Device
}
