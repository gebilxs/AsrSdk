@echo off
del c_demo.exe
gcc -m32 ./c_demo.c -o c_demo.exe -I. -L. -lsoe
c_demo.exe
if "%1" == "1" (
 cd ..
)