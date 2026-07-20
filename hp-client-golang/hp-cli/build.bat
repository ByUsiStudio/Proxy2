@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ../../target/proxy2.exe main.go

SET CGO_ENABLED=0
SET GOOS=android
SET GOARCH=arm64
go build -o ../../target/proxy2-android-arm64 main.go

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o ../../target/proxy2-386.exe main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -o ../../target/proxy2-386 main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ../../target/hp-lite-amd64 main.go

set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -o ../../target/proxy2-apple-amd64 main.go

set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=arm64
go build -o ../../target/proxy2-apple-arm64 main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -o ../../target/proxy2-arm64 main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o ../../target/proxy2-armv7 main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=mipsle
set GOMIPS=softfloat
go build -o ../../target/proxy2-mipsle main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=mips64le
set GOMIPS=softfloat
go build -o ../../target/hp-lite-mips64le main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=mips
set GOMIPS=softfloat
go build -o ../../target/proxy2-mips main.go

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=mips64
set GOMIPS=softfloat
go build -o ../../target/proxy2-mips64 main.go
