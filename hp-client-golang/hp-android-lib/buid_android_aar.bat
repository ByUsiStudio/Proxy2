@echo off
SET ANDROID_HOME=C:\Users\17782\AppData\Local\Android\Sdk
SET ANDROID_NDK_HOME=C:\Users\17782\AppData\Local\Android\Sdk\ndk\30.0.15729638
go install golang.org/x/mobile/cmd/gomobile
go get golang.org/x/mobile/bind
gomobile init
gomobile bind -target=android -androidapi 21
