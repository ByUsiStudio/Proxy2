export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o ../../target/proxy2.exe main.go

export CGO_ENABLED=0
export GOOS=android
export GOARCH=arm64
go build -o ../../target/proxy2-android-arm64 main.go

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=386
go build -o ../../target/proxy2-386.exe main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=386
go build -o ../../target/proxy2-386 main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o ../../target/proxy2-amd64 main.go

export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=amd64
go build -o ../../target/proxy2-apple-amd64 main.go

export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=arm64
go build -o ../../target/proxy2-apple-arm64 main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64
go build -o ../../target/proxy2-arm64 main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
export GOARM=7
go build -o ../../target/proxy2-armv7 main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=mipsle
export GOMIPS=softfloat
go build -o ../../target/proxy2-mipsle main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=mips64le
export GOMIPS=softfloat
go build -o ../../target/proxy2-mips64le main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=mips
export GOMIPS=softfloat
go build -o ../../target/proxy2-mips main.go

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=mips64
export GOMIPS=softfloat
go build -o ../../target/proxy2-mips64 main.go
