@echo off
REM Photon CLI wrapper for Windows
REM This script handles calling the downloaded Go binary

setlocal enabledelayedexpansion

set SCRIPT_DIR=%~dp0
set BINARY=%SCRIPT_DIR%photon-cli.exe

if not exist "%BINARY%" (
  echo ❌ Photon CLI binary not found. Please reinstall: npm install -g photon-cli
  exit /b 1
)

REM Execute the binary with all passed arguments
"%BINARY%" %*
