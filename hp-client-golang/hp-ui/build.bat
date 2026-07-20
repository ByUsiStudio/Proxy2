@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
fyne package -os windows -icon Icon.png -name proxy2.exe main.go
