@echo off
echo Building Windows Runner with Version Info...

REM Compile the resource file
windres -i version.rc -o version.syso

REM Build the executable with version info
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-H=windowsgui -s -w" -o runner.exe runner.go

REM Clean up
if exist version.syso del version.syso

echo Build complete: runner.exe