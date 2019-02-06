#!/bin/sh
set -e

# first arg is `-f` or `--some-option`
if [ "${1#-}" != "$1" ]; then
    set -- dbp "$@"
fi

if dbp "$1" --help 2>&1 >/dev/null | grep "help requested" > /dev/null 2>&1; then
    set -- dbp "$@"
fi

exec "$@"
