package bluetooth

import (
	"fmt"
	"sync"

	"github.com/KoruptTinker/audio-sync/internal/config"
	"tinygo.org/x/bluetooth"
)

func NewScanner(cfg config.BluetoothConfig) *Scanner {
	return &Scanner{
		cfg:     cfg,
		adapter: bluetooth.DefaultAdapter,
	}
}

func (s *Scanner) Adapter() *bluetooth.Adapter {
	return s.adapter
}

func (s *Scanner) Enable() error {
	if err := s.adapter.Enable(); err != nil {
		return fmt.Errorf("enable bluetooth adapter: %w", err)
	}
	return nil
}

func (s *Scanner) ScanForDevices() ([]bluetooth.Address, error) {
	devices := make(map[string]bluetooth.Address)
	var mu sync.Mutex

	fmt.Printf("Scanning for %d device(s) named %q...\n", s.cfg.DeviceCount, s.cfg.DeviceName)

	err := s.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() != s.cfg.DeviceName {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		addrStr := result.Address.String()
		if _, exists := devices[addrStr]; !exists {
			devices[addrStr] = result.Address
			fmt.Printf("Found device: %s\n", addrStr)
		}

		if len(devices) >= s.cfg.DeviceCount {
			if err := adapter.StopScan(); err != nil {
				fmt.Printf("Warning: failed to stop scan: %v\n", err)
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("scan for devices: %w", err)
	}

	addresses := make([]bluetooth.Address, 0, len(devices))
	for _, addr := range devices {
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

func Connect(adapter *bluetooth.Adapter, addr bluetooth.Address, cfg config.BluetoothConfig) (*Connection, error) {
	device, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, fmt.Errorf("connect to device: %w", err)
	}

	conn := &Connection{
		device: device,
		cfg:    cfg,
	}

	if err := conn.discoverCharacteristic(); err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Connection) discoverCharacteristic() error {
	services, err := c.device.DiscoverServices(nil)
	if err != nil {
		return fmt.Errorf("discover services: %w", err)
	}

	for _, service := range services {
		characs, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			return fmt.Errorf("discover characteristics: %w", err)
		}

		for _, charac := range characs {
			if charac.String() == c.cfg.CharacteristicUUID {
				c.characteristic = charac
				return nil
			}
		}
	}

	return fmt.Errorf("characteristic %s not found", c.cfg.CharacteristicUUID)
}

func (c *Connection) WriteWithoutResponse(data []byte) error {
	_, err := c.characteristic.WriteWithoutResponse(data)
	return err
}
