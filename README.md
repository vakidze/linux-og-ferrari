# Linux-og-ferari (GUI Serial Logger)

This repo builds a GTK-backed Fyne GUI app on Linux to read serial data and save CSVs.
Module path is `github.com/vakidze/Linux-og-ferari` — make sure repository uses this name.

Build locally:
```
sudo apt update
sudo apt install -y build-essential libgl1-mesa-dev xorg-dev libxkbcommon-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev pkg-config
go mod tidy
go build -o serial-logger ./cmd/app
./serial-logger
```

Push to GitHub and Actions will build a Linux binary and upload as an artifact.
