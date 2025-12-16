@echo off
REM Bitrim wrapper for Windows
REM This script handles calling the downloaded Go binary

setlocal enabledelayedexpansion

set SCRIPT_DIR=%~dp0
set BINARY=%SCRIPT_DIR%bitrim.exe

if not exist "%BINARY%" (
  echo ❌ Bitrim binary not found. Please reinstall: npm install -g bitrim
  exit /b 1
)

REM Execute the binary with all passed arguments
"%BINARY%" %*
