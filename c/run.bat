@echo off
del c_demo.exe
gcc ./c_demo.c -o c_demo.exe
c_demo.exe
if "%1" == "1" (
 cd ..
)