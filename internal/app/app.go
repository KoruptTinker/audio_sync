package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/KoruptTinker/audio-sync/internal/audio"
	"github.com/KoruptTinker/audio-sync/internal/bluetooth"
	"github.com/KoruptTinker/audio-sync/internal/config"
	"github.com/KoruptTinker/audio-sync/internal/led"
	bt "tinygo.org/x/bluetooth"
)

func Run() error {
	cfg := config.Default()

	dataChannels := createDataChannels(cfg.Bluetooth.DeviceCount)

	startAudioCapture(cfg.Audio, dataChannels)

	scanner := bluetooth.NewScanner(cfg.Bluetooth)
	if err := scanner.Enable(); err != nil {
		return fmt.Errorf("enable bluetooth: %w", err)
	}

	addresses, err := scanner.ScanForDevices()
	if err != nil {
		return fmt.Errorf("scan for devices: %w", err)
	}

	if len(addresses) == 0 {
		return fmt.Errorf("no devices found")
	}

	fmt.Printf("Connecting to %d device(s)...\n", len(addresses))

	return runLEDControllers(scanner, addresses, dataChannels, cfg)
}

func createDataChannels(count int) []chan float64 {
	channels := make([]chan float64, count)
	for i := range channels {
		channels[i] = make(chan float64, 1)
	}
	return channels
}

func startAudioCapture(cfg config.AudioConfig, dataChannels []chan float64) {
	handler := audio.NewRMSHandler(cfg, dataChannels)
	go func() {
		if err := handler.Start(); err != nil {
			log.Printf("Audio error: %v", err)
		}
	}()
}

func runLEDControllers(scanner *bluetooth.Scanner, addresses []bt.Address, dataChannels []chan float64, cfg config.Config) error {
	var wg sync.WaitGroup
	manager := led.NewManager(cfg.LED, cfg.Bluetooth)

	for i, addr := range addresses {
		conn, err := bluetooth.Connect(scanner.Adapter(), addr, cfg.Bluetooth)
		if err != nil {
			log.Printf("Failed to connect to %s: %v", addr.String(), err)
			continue
		}

		wg.Add(1)
		go manager.StartController(conn, dataChannels[i], &wg)
	}

	wg.Wait()
	return nil
}
