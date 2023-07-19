@echo off

REM Get the path of the startup folder
for /f "tokens=2*" %%A in ('reg query "HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders" /v Startup ^| find "Startup"') do (
    set "STARTUP_FOLDER=%%B"
)

REM Create or overwrite the batch script in the startup folder
echo @echo off > "%STARTUP_FOLDER%\magang-absen.bat"
echo cd %~dp0 >> "%STARTUP_FOLDER%\magang-absen.bat"
echo start magang-absen-otomatis.exe >> "%STARTUP_FOLDER%\magang-absen.bat"

REM Execute the magang-absen-otomatis.exe in the current directory
start magang-absen-otomatis.exe