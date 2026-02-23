@echo off
cd /d "%~dp0"
echo ========================================
echo  CyberBasic 3D Physics Demo
echo ========================================
echo.
if not exist cyberbasic.exe (
    echo cyberbasic.exe not found. Build first: go build -o cyberbasic.exe .
    pause
    exit /b 1
)
echo Starting... (window may take a moment)
echo If NO WINDOW appears: update graphics/OpenGL drivers, or run:
echo   go run test_raylib_window.go
echo to test raylib on this PC.
echo.
cyberbasic.exe examples\run_3d_physics_demo.bas
echo.
echo Exit code: %ERRORLEVEL%
pause
