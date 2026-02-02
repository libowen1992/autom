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


# 下载安装包
echo "-- 开始下载安装包${APP_PKG_URL}"
mkdir -p ${PKG_DOWNLOAD_DIR}
cd ${PKG_DOWNLOAD_DIR} && wget -q -O ${APP_PKG_NAME} ${APP_PKG_URL}

NODES={{ .Nodes }}
for node in ${NODES[@]}
    do
      echo "-- ${node} 复制安装包到远程临时目录 ${REMOTE_TMP_DIR}"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "mkdir -p ${REMOTE_TMP_DIR}"
      scp -q -o StrictHostKeyChecking=no -q  -P ${PORT} ./${APP_PKG_NAME} ${USER}@${node}:${REMOTE_TMP_DIR}
      echo "-- ${node} 备份服务"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "mkdir -p ${REMOTE_BACKUP_DIR}; cp -a  ${APP_PKG_PATH}/${APP_NAME} ${REMOTE_BACKUP_DIR}"
      echo "-- ${node} 替换新的包"
      ssh -q -o StrictHostKeyChecking=no -p ${PORT} ${USER}@${node} "cd ${APP_PKG_PATH} && rm -fr ${APP_NAME} && mkdir -p ${APP_NAME} && cd ${APP_NAME} && mv ${REMOTE_TMP_DIR}/${APP_PKG_NAME} . && tar xf ${APP_PKG_NAME} && rm -f ${APP_PKG_NAME}"
done