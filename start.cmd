: #! /bin/env sh # Heheheheh this is evil, it throws a soft error but works on both linux and windows 
:<<"::CMDCODE"
GOTO :WINDOWSSCRIPT
::CMDCODE

eval "$(cat /etc/*-release | awk '$0="export "$0')"
echo "[96m -- MCWrapper v0.1-alpha CLI $PRETTY_NAME -- [0m"
go run ./src/*.go

exit $?

:WINDOWSSCRIPT
for /F "tokens=*" %n IN ('go version') DO @(ECHO [46m -- MCWrapper v0.1-alpha CLI Windows %s -- [0m)
go run src/*.go