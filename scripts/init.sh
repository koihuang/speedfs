#! /bin/bash
set -e

base_dir=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
export SPEEDFS_PATH=${base_dir}

./speedfs init