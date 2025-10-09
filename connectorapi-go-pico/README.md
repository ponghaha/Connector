# connectorapi-go

A Go project built using the [Gin Web Framework](https://github.com/gin-gonic/gin).  

## 🚧 Build and Deploy
 
Before running in production, you should compile the project into a binary and move it to your desired folder.
 
### ✅ Step 1: Build the project
```bash
go build -o connector-api ./cmd/server/main.go
```
```bash
or
GOOS=linux GOARCH=amd64 go build -o connector-api ./cmd/server/main.go
```
 
### ✅ Step 2: Create deployment directory and move files for VM
```bash
mkdir -p /opt/connector-api
mkdir -p /opt/connector-api/configs
mkdir -p /opt/connector-api/elk/log

put connector-api to /connector-api
```

## 🚀 How to Run

This project is designed to run on a virtual machine (VM). You can run it directly using:

```bash
go run main.go [env]

For background execution (recommended for production), you may use:

<!--
If the binary cannot be executed (e.g., "Permission denied"), run:
chmod +x ./connector-api
This adds executable permission to the binary.
-->

1.create connector-api.service and copy below
 
[Unit]
Description=Sidecar Application (GoLang)
After=network.target

[Service]
User=root
Group=root
WorkingDirectory=/opt/connector-api
ExecStart=/bin/bash -c '/opt/connector-api/connector-api >> /opt/connector-api/log/applicationlog.txt 2>&1'
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
 
2. move to /etc/systemd/system/connector-api.service
3. sudo systemctl daemon-reload
4. sudo systemctl start connector-api.service
5. systemctl status connector-api.service -l or lsof -i -P -n | grep LISTEN


⚙️ Configuration
All environment-specific settings are placed under the /config directory with the filename format:
{config}.yaml

You must pass the environment name as a CLI parameter when running the program.


🧩 Dependencies
Make sure to install Go modules before running the project:
go mod tidy


👨‍💻 Author
SYE Section
Mr. Akkharasarans


📁 Project Structure
.
├───cmd
│   └───server
│       └───main.go
├───configs
│   └───config.yaml
├───docs
├───elk
│   └───log
├───internal
│   ├───adapter
│   │   ├───client
│   │   ├───handler
│   │   │   └───api
│   │   └───utils
│   └───core
│       ├───domain
│       └───service
│           └───format
└───pkg
    ├───config
    ├───error
    ├───format
    ├───logger
    └───metrics
