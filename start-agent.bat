@echo off
cd /d %~dp0
set WOODPECKER_SERVER=localhost:9000
set WOODPECKER_AGENT_SECRET=mysecret123
set WOODPECKER_BACKEND=local
set WOODPECKER_MAX_WORKFLOWS=4
set WOODPECKER_LOG_LEVEL=info
set WOODPECKER_GRPC_SECURE=false
dist\woodpecker-agent.exe
