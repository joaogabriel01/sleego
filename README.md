# Sleego

![Sleego](docs/images/sleego_doc.png)

Sleego is a small Go-based tool that enforces **time-based rules** on computer usage.

It works by:

* monitoring running processes
* applying allowed time windows
* terminating processes that run outside their configured schedule
* optionally triggering a system shutdown at a fixed time

Sleego is **configuration-driven** and does not provide a graphical interface in this repository.

---

## Overview

Sleego continuously monitors system processes based on a JSON configuration file.

For each configured application (or logical rule), Sleego:

* checks whether it is running
* verifies if the current time is within the allowed interval
* terminates the process if it is outside that interval

Additionally, Sleego can schedule a system shutdown and emit warnings before it happens.

The project is designed to be predictable and non-interactive at runtime.
All decisions are made **beforehand**, in the configuration file.

---

## Features

* **Process scheduling**

  * Define allowed time windows for applications
  * Processes running outside their window are terminated

* **Scheduled system shutdown**

  * Define a fixed shutdown time
  * Receive advance warnings
  * Shutdown happens automatically

* **Categories (logical rules)**

  * Applications can be grouped under logical names
  * Categories behave like application rules

* **Notifications**

  * Alerts before shutdown
  * Alerts when processes are terminated

---

## Configuration

Sleego is fully driven by a `config.json` file.

### Example

```json
{
  "apps": [
    {
      "name": "browser.exe",
      "allowed_from": "09:00",
      "allowed_to": "18:00"
    },
    {
      "name": "games",
      "allowed_from": "20:00",
      "allowed_to": "23:30"
    }
  ],
  "shutdown": "23:59",
  "categories": {
    "games": ["steam.exe", "game.exe"]
  }
}
```

### Fields

* **apps**

  * `name`: process name or logical category name
  * `allowed_from`: start time (HH:MM)
  * `allowed_to`: end time (HH:MM)

* **shutdown**

  * Time when the system should shut down (HH:MM)

* **categories**

  * Map of logical names to process names
  * Categories can be referenced in `apps` like regular applications

---

## Usage (CLI)

Sleego is intended to be run as a long-lived process.

### Basic execution

```bash
make cli
./sleego
```

Once started, Sleego:

* loads the configuration
* starts monitoring processes
* applies shutdown rules
* runs indefinitely

Any change to the configuration file requires **restarting the process**.

---

## Notifications

Sleego emits notifications for:

* Process termination events
* Upcoming shutdown warnings (for example, 10 minutes before)

Notifications are informational only and do not require user interaction.

---

## Security and considerations

* **Forceful termination**

  * Processes are killed immediately if they violate the schedule
  * Double-check configuration to avoid unintended data loss

* **Shutdown**

  * The system will shut down at the configured time
  * Save your work beforehand

* **Permissions**

  * Sleego must run with sufficient privileges to:

    * terminate processes
    * trigger system shutdown

---

## Intended usage

Sleego is intentionally simple and opinionated.

It is not:

* a productivity tracker
* a parental control system
* an analytics or reporting tool

It exists to enforce **predefined boundaries** with minimal runtime interaction.

---

## Related projects

* **Sleego UI**
  A separate repository provides a desktop UI built on top of this core.

  [https://github.com/joaogabriel01/sleego-ui](https://github.com/joaogabriel01/sleego-ui)

---

## License

This project is licensed under the MIT License.

---

### Notes on scope

This repository contains **only the core engine**.
Any user interface, installer, or system integration lives outside of this project.

