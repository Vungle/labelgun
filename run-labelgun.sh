#!/bin/sh
set -e

if [ "${LABELGUN_SUPPRESS_LOG}" = "true" ]; then
  labelgun >/dev/null 2>&1
  exit $?
fi

labelgun
