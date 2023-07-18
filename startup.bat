@echo off
set "TARGET_FILE=magang-absen-otomatis.exe"

REM Get the path of the startup folder
for /f "tokens=2*" %%A in ('reg query "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders" /v Startup ^| find "Startup"') do (
    set "STARTUP_FOLDER=%%B"
)

REM Create a symbolic link in the startup folder
set "LINK_FILE=%STARTUP_FOLDER%\%TARGET_FILE%"
set "TARGET_PATH=%CD%\%TARGET_FILE%"
set "CONFIG_FILE=%CD%\config.yaml"
set "IMAGE_FILE=%CD%\absen.png"

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

REM Create a symbolic link for config.yaml
set "CONFIG_LINK=%STARTUP_FOLDER%\config.yaml"
echo "%CONFIG_LINK%"
if exist "%CONFIG_LINK%" (
    echo Symbolic link already exists for config.yaml.
) else (
    echo Creating symbolic link for config.yaml...
    mklink "%CONFIG_LINK%" "%CONFIG_FILE%"
    if %errorlevel% equ 0 (
        echo Symbolic link created for config.yaml.
    ) else (
        echo Failed to create symbolic link for config.yaml.
    )
)

REM Create a symbolic link for absen.png
set "IMAGE_LINK=%STARTUP_FOLDER%\absen.png"
echo "%IMAGE_LINK%"
if exist "%IMAGE_LINK%" (
    echo Symbolic link already exists for absen.png.
) else (
    echo Creating symbolic link for absen.png...
    mklink "%IMAGE_LINK%" "%IMAGE_FILE%"
    if %errorlevel% equ 0 (
        echo Symbolic link created for absen.png.
    ) else (
        echo Failed to create symbolic link for absen.png.
    )
)

pause
