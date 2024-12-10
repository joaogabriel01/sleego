# Sleego

**Sleego** is a Go application designed to monitor and control the execution of processes based on specified schedules. It allows you to configure time restrictions for applications, forcefully terminating those running outside their permitted hours. Additionally, Sleego supports scheduled system shutdowns.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Screenshots](#screenshots)
- [Configuration Structure](#configuration-structure)
- [Installation and Execution](#installation-and-execution)
  - [For End Users (Windows 64-bit Installer)](#for-end-users-windows-64-bit-installer)
  - [For Advanced Users (Run from Source)](#for-advanced-users-run-from-source)
- [Usage](#usage)
  - [Process Monitoring](#process-monitoring)
  - [Scheduled Shutdown](#scheduled-shutdown)
- [Notifications](#notifications)
- [Security and Considerations](#security-and-considerations)
- [Contributions](#contributions)
  - [Todo](#todo)
- [License](#license)

## Overview

Sleego monitors running processes based on a JSON configuration file that lists processes and their allowed execution times. It periodically checks for running processes and immediately terminates any outside their permitted schedule.

- **GUI Version:** Modify configuration directly in the interface. After editing, click "Save and Run" to apply changes without restarting.
- **CLI Version:** Requires restarting the application after editing the configuration file.

Both versions support scheduled system shutdowns and provide system notifications about terminated processes and upcoming shutdowns.

## Features

- **Process Monitoring:** Terminates processes running outside allowed time frames.
- **Scheduled Shutdown:** Shuts down the computer at a specified time, with prior notifications.
- **System Notifications:** Alerts when processes are terminated and before system shutdown.
- **GUI with System Tray:** Easy access and user-friendly interface.

## Screenshots

![Sleego Main Interface](images/sleego_main_interface.png)  
*Figure 1: Sleego's main GUI interface.*

## Configuration Structure

Your `config.json` should follow this format:

```json
{
  "apps": [
    {
      "name": "example1.exe",
      "allowed_from": "09:00",
      "allowed_to": "10:00"
    },
    {
      "name": "example2.exe",
      "allowed_from": "14:00",
      "allowed_to": "23:30"
    }
  ],
  "shutdown": "23:59"
}
```

- **name**: The process name (e.g., `app1.exe`).
- **allowed_from**: Allowed start time (HH:MM).
- **allowed_to**: Allowed end time (HH:MM).
- **shutdown**: Scheduled shutdown time.

## Installation and Execution

### For End Users (Windows 64-bit Installer)

If you’re on Windows 64-bit, use the provided installer for a hassle-free setup:

1. **Run the Installer:**  
   Double-click the installer (setup/) and follow the on-screen instructions.

2. **Launch Sleego:**  
   Once installed, you can run Sleego from the Start Menu or by navigating to its installation folder and double-clicking the executable.

### For Advanced Users (Run from Source)

If you prefer to run Sleego from source (for development, customization, or non-Windows platforms):

**Prerequisites:**
- **Go Installed (v1.18+ recommended)**
- **Ensure `images` and `config.json` Exist** in the working directory.
- **GCC for GUI Builds** (if you want the GUI version):
  - **Windows:** Install [MinGW](https://www.mingw-w64.org/downloads/)
  - **macOS:**  
    ```bash
    xcode-select --install
    ```
  - **Linux (Debian-based):**  
    ```bash
    sudo apt-get install build-essential
    ```

**Building Steps:**
1. **Clone the repository**:
    ```bash
    git clone https://github.com/joaogabriel01/sleego.git
    cd sleego
    ```
2. **Install dependencies**:
    ```bash
    go mod tidy
    ```
3. **Build using Go**:
    - **CLI Version**:
      ```bash
      go build -o sleego_cli ./cmd/cli/main.go
      ```
    - **GUI Version**:
      ```bash
      go build -o sleego_gui ./cmd/gui/main.go
      ```

**Using the Makefile:**  
A Makefile is included to streamline the build process. For example:
```bash
make cli
make linux_gui
make windows_gui
```
`make clean` removes build artifacts.

## Usage

### Process Monitoring

When Sleego starts, it loads the configuration and monitors processes. Processes running outside allowed times are terminated immediately.

**Example:**
```
Starting process policy with config: [{name:app1.exe, allowed_from:08:00, allowed_to:18:00}, {name:app2.exe, allowed_from:09:00, allowed_to:23:20}] from path: ./config.json
```

### Scheduled Shutdown

- **GUI:**
  1. Open the Sleego GUI.
  2. Enter the shutdown time (HH:MM).
  3. Click "Run".

- **CLI:**
  1. Edit the `shutdown` field in `config.json`.
  2. Restart the CLI application.

Sleego will display notifications as the shutdown time approaches and before initiating it.

## Notifications

- **Process Termination Alerts:**  
  Alerts when a process is terminated due to schedule restrictions.

- **Shutdown Warnings:**  
  Alerts before the scheduled shutdown (e.g., 10, 3, and 1 minute warnings).

## Security and Considerations

- **Forceful Termination:**  
  Processes are closed without warning. Double-check your configuration.
- **Scheduled Shutdown:**  
  The system will shut down at the specified time—save your work in advance.
- **Permissions:**  
  Ensure you have permissions to terminate listed processes and to shut down the system.
- **System Tray Access:**  
  Sleego runs in the background and is accessible from the system tray.

## Contributions

Contributions are welcome. Feel free to open an issue or submit a pull request.

### Todo

- **Application Groups**: Ability to create groups of applications with the same schedule.
- **Process Selection UI**: Ability to view running processes to facilitate adding processes to be monitored.
- **Enhanced GUI Visualization**: Add options for customizing themes (light/dark mode) and resizing the interface for better usability and readability.

## License

This project is licensed under the MIT License.