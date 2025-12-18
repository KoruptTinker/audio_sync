package main

import (
	"fmt"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

func main() {
	wg := sync.WaitGroup{}
	adapter := bluetooth.DefaultAdapter
	connectionHandles := map[string]bluetooth.Address{}
	must("Enable Adapter", adapter.Enable())
	fmt.Println("Scanning for devices...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() == "GATT--DEMO" {
			_, exists := connectionHandles[result.Address.String()]
			if !exists {
				connectionHandles[result.Address.String()] = result.Address
				fmt.Printf("Found device: %s \n", result.Address.String())
			}
		}
		if len(connectionHandles) == 2 {
			stopErr := adapter.StopScan()
			must("Stop Scan", stopErr)
		}
	})

	must("Scan devices", err)
	fmt.Println("Connecting to devices...")
	wg.Add(2)
	for _, addr := range connectionHandles {
		go handleConnection(addr, adapter, &wg)
	}
	wg.Wait()
}

func must(actionName string, err error) {
	if err != nil {
		panic(fmt.Sprintf("Failed to complete action: %s, ERROR: %s", actionName, err.Error()))
	}
}

func handleConnection(addr bluetooth.Address, adapter *bluetooth.Adapter, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Handling connection")
	device, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
	must("Connect to device", err)
	var characteristic bluetooth.DeviceCharacteristic
	services, err := device.DiscoverServices(nil)
	must("Discover Services", err)
	for _, service := range services {
		characs, err := service.DiscoverCharacteristics(nil)
		must("Discover Characteristics", err)
		for _, charac := range characs {
			if charac.String() == "0000fff3-0000-1000-8000-00805f9b34fb" {
				characteristic = charac
			}
		}
	}

	offPayload := []byte{0xBC, 0x01, 0x01, 0x00, 0x55}
	onPayload := []byte{0xBC, 0x01, 0x01, 0x01, 0x55}
	characteristic.WriteWithoutResponse(offPayload)
	time.Sleep(500 * time.Millisecond)
	characteristic.WriteWithoutResponse(onPayload)
	time.Sleep(30 * time.Millisecond)
}
