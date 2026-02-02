#!/usr/bin/env bash

num=`netstat -tunlp|grep 10099|grep -v grep|wc -l`
if [ $num -eq 1 ];then
    echo "autom web has running"
  else
    ./gotty --config ./.gotty ./autom  >> ./gotty.log 2>&1 &
fi

