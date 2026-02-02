#!/usr/bin/env bash

USER={{ .User }}
PORT={{ .Port }}
REMOTE_BACKUP_DIR={{ .RemoteBackupDir }}
APP_PKG_NAME={{ .AppPkgName }}
APP_PKG_PATH={{ .AppPkgPath }}
STOP_CMD="{{ .StopCmd }}"
START_CMD="{{ .StartCmd }}"

NODES={{ .Nodes }}
for node in ${NODES[@]}
    do
      echo "-- ${node} 停止服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "${STOP_CMD}"
      echo "-- ${node} 替换新的包"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "cd ${APP_PKG_PATH} && rm -fr ${APP_PKG_NAME} && cp -a ${REMOTE_BACKUP_DIR}/${APP_PKG_NAME} ."
      echo "-- ${node} 启动服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "source /etc/profile && ${START_CMD}"
done