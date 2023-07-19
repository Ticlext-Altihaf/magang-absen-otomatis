#!/bin/bash

# Create a systemd service file
sudo tee /etc/systemd/system/magang-absen.service > /dev/null <<EOF
[Unit]
Description=Magang Absen Service
After=network.target

[Service]
Type=simple
ExecStart=$(pwd)/magang-absen-otomatis.exe

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd daemon
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable magang-absen.service

# Start the service
sudo systemctl start magang-absen.service