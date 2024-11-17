# Sleego

**Sleego** is a Go application designed to monitor and control the execution of processes based on specified schedules. It allows you to configure time restrictions for applications, forcefully terminating those running outside their permitted hours.

## Overview

Sleego reads a JSON configuration file that lists the processes to monitor and the allowed execution times for each. Based on this information, Sleego periodically checks the running processes and immediately terminates any that are outside their specified schedule.

> **Warning**: Sleego forcefully closes processes that are outside the allowed schedule without any prior warning. Therefore, avoid adding critical processes without confirming the time settings.

## Configuration Structure

The configuration JSON file should follow this format:

```json
[
    {
        "name": "app1.exe",
        "allowed_from": "08:00",
        "allowed_to": "18:00"
    },
    {
        "name": "app2.exe",
        "allowed_from": "09:00",
        "allowed_to": "23:20"
    }
]
```

- **name**: The name of the process (e.g., `app1.exe`).
- **allowed_from**: The allowed start time in HH:MM format.
- **allowed_to**: The allowed end time in HH:MM format.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/joaogabriel01/sleego.git
    cd sleego
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. **For GUI Execution**:
    - Install GCC (required for the GUI version).
    - On Windows, install MinGW.
    - On macOS, install Xcode Command Line Tools:
        ```bash
        xcode-select --install
        ```
    - On Linux based in Debian:
        ```bash
        sudo apt-get install build-essential
        ```

## Execution Methods

Sleego can be executed in two ways:

### CLI Execution

Run the command-line interface version located at `cmd/cli/main.go`:

```bash
go run ./cmd/cli/main.go -config="./config.json"
```
or
```bash
go build ./cmd/cli/main.go
./main or ./main.exe
```

### GUI Execution

Run the graphical user interface version located at `cmd/gui/main.go`:

```bash
go run ./cmd/gui/main.go
```
or
```bash
go build ./cmd/gui/main.go
./main or ./main.exe
```

**Note**: Ensure GCC is installed before running the GUI version.

## Execution Example

When Sleego starts, it loads the provided configuration file and begins monitoring processes. Example output:

```
Starting process policy with config: [{name:app1.exe, allowed_from:08:00, allowed_to:18:00}, {name:app2.exe, allowed_from:09:00, allowed_to:23:20}] of path: ./config.json
```

## Security and Considerations

- **Processes are forcefully terminated**: Sleego terminates any processes that do not meet the allowed times without warning. Ensure the configuration is correct to prevent terminating essential processes.
- **Permissions**: Make sure you have the necessary permissions to close the listed processes.

## Contributions

Contributions are welcome. Feel free to open an issue or submit a pull request.

### Todo

- Ability to create groups of applications with the same schedule.
- Ability to view running processes to facilitate adding processes to be monitored.

## License

This project is licensed under the MIT License.