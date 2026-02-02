#!/usr/bin/env bash

USER={{ .User }}
PORT={{ .Port }}
REMOTE_BACKUP_DIR={{ .RemoteBackupDir }}
APP_NAME={{ .AppName }}
APP_PKG_PATH={{ .AppPkgPath }}


NODES={{ .Nodes }}
for node in ${NODES[@]}
    do
      echo "-- ${node} 替换回滚包"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "cd ${APP_PKG_PATH} && rm -fr ${APP_NAME} && cp -a ${REMOTE_BACKUP_DIR}/${APP_NAME} ."
done