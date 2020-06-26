#!/bin/sh
set -e

if [ -z "${LABELGUN_ERR_THRESHOLD}" ]; then
  LABELGUN_ERR_THRESHOLD="WARNING"
fi

labelgun -stderrthreshold="${LABELGUN_ERR_THRESHOLD}"
