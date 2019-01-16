#!/bin/sh
set -e

INSTALLATION_PATH=/usr/local/bin/autoscaler-agent
UNIT_PATH=/etc/systemd/system/autoscaler-agent.service

wget -O $INSTALLATION_PATH $AGENT_URI
chmod +x $INSTALLATION_PATH

# FIXME verify checksum
# FIXME non-systemd systems
cat >$UNIT_PATH <<EOF
[Unit]
Description=CloudAutoScaler daemon
After=network.target
           
[Service]
User=nobody
ExecStart=$INSTALLATION_PATH
Environment=CAS_URI=$CAS_URI
Environment=CAS_TOKEN=$CAS_TOKEN
Restart=always
           
[Install]
WantedBy=multi-user.target
EOF
chmod 640 $UNIT_PATH
systemctl daemon-reload
systemctl start autoscaler-agent.service
systemctl enable autoscaler-agent.service
