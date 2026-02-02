#!/usr/bin/env bash

USER={{ .User }}
PORT={{ .Port }}
REMOTE_TMP_DIR={{ .RemoteTmpDir }}
REMOTE_BACKUP_DIR={{ .RemoteBackupDir }}
APP_NAME={{ .AppName }}
APP_PKG_NAME={{ .AppPkgName }}
APP_PKG_URL={{ .PkgUrl }}
APP_PKG_PATH={{ .AppPkgPath }}
PKG_DOWNLOAD_DIR={{ .PkgDownLoadDir }}
STOP_CMD="{{ .StopCmd }}"
START_CMD="{{ .StartCmd }}"

# 下载安装包
echo "-- 开始下载安装包${APP_PKG_URL}"
cd ${PKG_DOWNLOAD_DIR} && wget -q -O ${APP_PKG_NAME} ${APP_PKG_URL}

NODES={{ .Nodes }}
for node in ${NODES[@]}
    do
      echo "-- ${node} 复制安装包到远程临时目录 ${REMOTE_TMP_DIR}"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "mkdir -p ${REMOTE_TMP_DIR}"
      scp -q -o StrictHostKeyChecking=no -q  -P ${PORT} ./${APP_PKG_NAME} ${USER}@${node}:${REMOTE_TMP_DIR}
      echo "-- ${node} 备份服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "mkdir -p ${REMOTE_BACKUP_DIR}; cp -a ${APP_PKG_PATH}/${APP_PKG_NAME} ${REMOTE_BACKUP_DIR}"
      echo "-- ${node} 停止服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "${STOP_CMD}"
      echo "-- ${node} 替换新的包"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "cd ${APP_PKG_PATH} && rm -fr ${APP_NAME} ${APP_PKG_NAME} && mv ${REMOTE_TMP_DIR}/${APP_PKG_NAME} ."
      echo "-- ${node} 启动服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "source /etc/profile && ${START_CMD}"
done