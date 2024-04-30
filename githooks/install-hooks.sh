#!/usr/bin/env bash

set -e
set -u

pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg
