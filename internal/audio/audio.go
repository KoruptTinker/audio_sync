package audio

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"unsafe"

	"github.com/KoruptTinker/audio-sync/internal/config"
	"github.com/gen2brain/malgo"
)

func NewRMSHandler(cfg config.AudioConfig, channels []chan float64) *RMSHandler {
	return &RMSHandler{
		cfg:      cfg,
		channels: channels,
	}
}

func (h *RMSHandler) Start() error {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return fmt.Errorf("init audio context: %w", err)
	}
	h.ctx = ctx

	deviceID, err := h.findDevice()
	if err != nil {
		return err
	}

	var pinner runtime.Pinner
	pinner.Pin(deviceID)
	defer pinner.Unpin()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = h.cfg.Channels
	deviceConfig.SampleRate = h.cfg.SampleRate
	deviceConfig.Alsa.NoMMap = 1
	deviceConfig.Capture.DeviceID = unsafe.Pointer(deviceID)

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: h.onAudioData,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return fmt.Errorf("init audio device: %w", err)
	}
	h.device = device

	if err := device.Start(); err != nil {
		return fmt.Errorf("start audio device: %w", err)
	}

	fmt.Println("Audio capture started. Press Enter to stop.")
	fmt.Scanln()

	return nil
}

func (h *RMSHandler) Stop() {
	if h.device != nil {
		h.device.Uninit()
	}
	if h.ctx != nil {
		h.ctx.Free()
	}
}

func (h *RMSHandler) findDevice() (*malgo.DeviceID, error) {
	infos, err := h.ctx.Devices(malgo.Capture)
	if err != nil {
		return nil, fmt.Errorf("enumerate devices: %w", err)
	}

	deviceID := new(malgo.DeviceID)

	fmt.Println("Scanning audio devices...")
	for _, info := range infos {
		if strings.Contains(info.Name(), h.cfg.DeviceName) {
			fmt.Printf("-> Found: %s\n", info.Name())
			*deviceID = info.ID
			return deviceID, nil
		}
	}

	return nil, fmt.Errorf("audio device %q not found", h.cfg.DeviceName)
}

func (h *RMSHandler) onAudioData(output, input []byte, frameCount uint32) {
	if len(input) == 0 {
		return
	}

	rms := CalculateRMS(input)

	for _, ch := range h.channels {
		select {
		case ch <- rms:
		default:
		}
	}

	printBar(rms)
}

func CalculateRMS(input []byte) float64 {
	sampleCount := len(input) / 4
	if sampleCount == 0 {
		return 0
	}

	samples := unsafe.Slice((*float32)(unsafe.Pointer(&input[0])), sampleCount)

	var sumSquares float64
	for _, sample := range samples {
		sumSquares += float64(sample * sample)
	}

	return math.Sqrt(sumSquares / float64(sampleCount))
}

func printBar(val float64) {
	scaled := int(val * 50 * 2.0)
	if scaled > 50 {
		scaled = 50
	}
	fmt.Printf("\rRMS: [%-50s] %.3f", strings.Repeat("|", scaled), val)
}
