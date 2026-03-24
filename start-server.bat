@echo off
cd /d %~dp0
set PATH=C:\Users\XOS\scoop\apps\gcc\15.2.0\bin;%PATH%
set WOODPECKER_OPEN=true
set WOODPECKER_ADMIN=admin
set WOODPECKER_HOST=https://nonrenouncing-nonsalably-shae.ngrok-free.dev
set WOODPECKER_AGENT_SECRET=mysecret123
set WOODPECKER_DATABASE_DRIVER=sqlite3
set WOODPECKER_DATABASE_DATASOURCE=%~dp0woodpecker.sqlite
set WOODPECKER_LOG_LEVEL=info
set WOODPECKER_GITHUB=true
set WOODPECKER_GITHUB_CLIENT=Ov23liZo1mvEjmkzlzea
set WOODPECKER_GITHUB_SECRET=7f98c01353aa37e09e126bf3dff4eb56f25d360f
dist\woodpecker-server.exe
