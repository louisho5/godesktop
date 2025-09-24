@echo off
echo Building Windows Runner with Version Info...

REM Compile the resource file
windres -o version.syso version.rc

REM Build the executable with version info
go build -ldflags "-H=windowsgui -s -w" -o runner.exe runner.go

REM Clean up
del version.syso

echo Build complete: runner.exe