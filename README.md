# CrOS/FW Chassis Monitor
> should be compatible with crOS laptops, tested on framework laptop - basically just shuts off the laptop instantly if the chassis is opened while the laptop is powered on, security/safety use case.

> ⚠️ chassis detection may cause rebooting unexpectedly, consider only starting the service before enabling it to check data is fetched correctly.

---
As a binary (required for service)

1. `go build -o chassis-monitor main.go`
2. `sudo mv chassis-monitor /usr/local/bin/chassis-monitor`
3. `sudo chmod 700 /usr/local/bin/chassis-monitor`
---
As a service

4. `sudo nano /etc/systemd/system/chassis-monitor.service`
```
[Unit]
Description=CrOS/FW Chassis Monitor
After=multi-user.target

[Service]
Type=simple
ExecStart=/usr/local/bin/chassis-monitor
Restart=always
RestartSec=3
User=root

[Install]
WantedBy=multi-user.target
```
5. `sudo systemctl daemon-reload`
6. `sudo systemctl enable --now chassis-monitor.service`
7. Confirm it's `active (running)` when running `sudo systemctl status chassis-monitor.service`
---