@echo off
set "TARGET_FILE=magang-absen-otomatis.exe"

REM Get the path of the startup folder
for /f "tokens=2*" %%A in ('reg query "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders" /v Startup ^| find "Startup"') do (
    set "STARTUP_FOLDER=%%B"
)

REM Create a symbolic link in the startup folder
set "LINK_FILE=%STARTUP_FOLDER%\%TARGET_FILE%"
set "TARGET_PATH=%CD%\%TARGET_FILE%"
echo "%LINK_FILE%"
REM Check if the symbolic link already exists
if exist "%LINK_FILE%" (
    echo Symbolic link already exists in the startup folder.
) else (
    echo Creating symbolic link in the startup folder...
    mklink "%LINK_FILE%" "%TARGET_PATH%"
    if %errorlevel% equ 0 (
        echo Symbolic link created in the startup folder.
    ) else (
        echo Failed to create symbolic link in the startup folder.
    )
)

pause
