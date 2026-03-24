@echo off
cd /d %~dp0
start start-web.bat
ping 127.0.0.1 -n 4 > nul
start start-server.bat
ping 127.0.0.1 -n 11 > nul
start start-agent.bat
pause
