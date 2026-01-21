# Audio Sync

This project synchronizes LED brightness with audio input in real-time. It captures system audio (or microphone input) and transmits brightness values to connected Bluetooth Low Energy (BLE) devices based on the audio's RMS (Root Mean Square) amplitude.

## Prerequisites

This application is currently designed for **macOS** systems.

* **macOS**: The project relies on macOS-specific audio routing via BlackHole.
* **Go**: You need the Go programming language installed (version 1.23 or later is recommended).
* **BlackHole**: A virtual audio driver is required to route system audio to this application.
  * Install via Homebrew: `brew install blackhole-2ch` (or 16ch/64ch)
  * Or download from the [BlackHole GitHub repository](https://github.com/ExistentialAudio/BlackHole).

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/KoruptTinker/audio-sync.git
    cd audio-sync
    ```

2. Download dependencies:

    ```bash
    go mod download
    ```

3. Build the application:

    ```bash
    go build -o audio-sync cmd/audio-sync/main.go
    ```

## Configuration

1. Create a configuration file from the example:

    ```bash
    cp config/config.example.yaml config.yaml
    ```

2. Edit `config.yaml` to match your setup:
    * **audio.devicename**: The name of your audio input device (e.g., "BlackHole 2ch").
    * **bluetooth.devicename**: The name of the BLE devices to connect to (e.g., "GATT--DEMO").
    * **bluetooth.devicecount**: The number of devices to wait for before starting.
    * **bluetooth.characteristicuuid**: The UUID of the characteristic to write brightness values to.
    * **led**: Adjust settings like `minbrightness`, `maxbrightness`, `gain`, and `threshold` to tune the LED response.

## Usage

1. **Configure Audio Output**:
    * Open **Audio MIDI Setup** on your Mac.
    * Create a **Multi-Output Device** that includes both your speakers/headphones and BlackHole.
    * Set this Multi-Output Device as your system output.
    * This allows you to hear the audio while sending it to BlackHole for capture.

2. **Run the Application**:

    ```bash
    ./audio-sync
    ```

3. The application will:
    * Scan for the configured Bluetooth devices.
    * Connect to them.
    * Start capturing audio from the configured input device.
    * Update the LED brightness on the connected devices in real-time.

4. Press **Enter** in the terminal to stop the audio capture and exit.

## Note on macOS Support

This tool explicitly relies on the **BlackHole** virtual audio driver for capturing system audio on macOS. While the core logic is written in Go, the audio routing workflow is specific to the macOS environment.
