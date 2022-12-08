@echo off
del libasr.dll
set GOARCH=386
set CGO_ENABLED=1
go build -ldflags "-s -w" -buildmode=c-shared -o libasr.dll
set GOARCH=amd64
IF %errorlevel% NEQ 0 GOTO ERROR
echo build dll success.
copy libasr.dll c
copy libasr.h c
GOTO END
:ERROR
    echo build dll failed.
:END

:: TODO 支持高阶题型