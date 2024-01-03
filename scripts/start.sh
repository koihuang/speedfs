#! /bin/bash
set -e

base_dir=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
export SPEEDFS_PATH=${base_dir}
cfg_file_path=${SPEEDFS_PATH}/config/config.json

function must_init(){
  if [ ! -f ${cfg_file_path} ]; then
    echo "speedfs not init, please execute './init.sh' first"
    exit
  fi
}
must_init

function start() {
    if [ -f ${SPEEDFS_PATH}/speedfs.pid ]; then
      pid=$(cat ${SPEEDFS_PATH}/speedfs.pid)
      if [ "${pid}" ]; then
        cmd_name=$(ps -p${pid} -o pid,comm | awk 'END{print $2}')
        if [[ "${cmd_name}" =~ "speedfs" ]]; then
          echo "speedfs has already be running, pid: ${pid}"
          return
        fi
      fi
    fi
    nohup ./speedfs start > /dev/null 2>&1 &
    echo "start speedfs, pid: $!"
}
start





