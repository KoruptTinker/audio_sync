package led

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/KoruptTinker/audio-sync/internal/bluetooth"
	"github.com/KoruptTinker/audio-sync/internal/config"
)

var (
	PayloadOff = []byte{0xBC, 0x01, 0x01, 0x00, 0x55}
	PayloadOn  = []byte{0xBC, 0x01, 0x01, 0x01, 0x55}
)

// NewController creates a new Controller with the given configuration.
func NewController(cfg config.LEDConfig, conn BLEWriter, dataChannel <-chan float64) *Controller {
	return &Controller{
		cfg:         cfg,
		conn:        conn,
		dataChannel: dataChannel,
	}
}

func (c *Controller) Initialize() error {
	if err := c.conn.WriteWithoutResponse(PayloadOff); err != nil {
		return fmt.Errorf("send off command: %w", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := c.conn.WriteWithoutResponse(PayloadOn); err != nil {
		return fmt.Errorf("send on command: %w", err)
	}
	time.Sleep(30 * time.Millisecond)

	return nil
}

func (c *Controller) Run() {
	ticker := time.NewTicker(c.cfg.UpdateInterval)
	defer ticker.Stop()

	payload := []byte{0xBC, 0x05, 0x06, 0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x55}
	var loudness float64

	for range ticker.C {
		select {
		case rms, ok := <-c.dataChannel:
			if !ok {
				return
			}
			loudness = rms
		default:
		}

		payload[3], payload[4] = c.CalculateBrightness(loudness)
		c.conn.WriteWithoutResponse(payload)
	}
}

func (c *Controller) CalculateBrightness(rms float64) (byte, byte) {
	if rms < c.cfg.Threshold {
		return 0, 100 // 0x0064 = minimum brightness
	}

	input := (rms - c.cfg.Threshold) * c.cfg.Gain

	// Clamp to [0, 1]
	input = math.Max(0, math.Min(1, input))

	// Apply power curve for deep contrast
	curve := math.Pow(input, c.cfg.CurveExponent)

	// Map to brightness range
	span := c.cfg.MaxBrightness - c.cfg.MinBrightness
	brightness := uint16(c.cfg.MinBrightness + (curve * span))

	return byte(brightness >> 8), byte(brightness & 0xFF)
}

func NewManager(cfg config.LEDConfig, bleCfg config.BluetoothConfig) *Manager {
	return &Manager{
		cfg:    cfg,
		bleCfg: bleCfg,
	}
}

func (m *Manager) StartController(conn *bluetooth.Connection, dataChannel <-chan float64, wg *sync.WaitGroup) {
	defer wg.Done()

	controller := NewController(m.cfg, conn, dataChannel)

	fmt.Println("Initializing LED...")
	if err := controller.Initialize(); err != nil {
		fmt.Printf("Failed to initialize LED: %v\n", err)
		return
	}

	fmt.Println("LED controller running...")
	controller.Run()
}
