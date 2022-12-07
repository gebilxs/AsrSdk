@echo off
del libsoe.dll
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags "-s -w" -buildmode=c-shared -o libsoe.dll
set GOARCH=amd64
IF %errorlevel% NEQ 0 GOTO ERROR
echo build dll success.
copy libsoe.dll c
copy libsoe.h c
cd c
run.bat 1
GOTO END
:ERROR
    echo build dll failed.
:END

:: TODO 支持高阶题型