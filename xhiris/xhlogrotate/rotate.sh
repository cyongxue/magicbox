#!/bin/sh

# logs目录不存在时就创建logs目录
if [ ! -d "logs/" ];then
  mkdir logs
  echo "create logs dir success"
else
  echo "logs dir already exists"
fi

# logs/rotate目录不存在时创建logs/rotate目录
if [ ! -d "logs/rotate/" ];then
  mkdir logs/rotate
  echo "create logs/rotate success"
else
  echo "logs/rotate dir already exists"
fi

modulepwd=$(pwd)

#### 以下3个参数，可以自定义调整，实现日志可配归档
logPrefix="server_"     # 日志文件前缀
rotateBeforHour=2       # 多长时间以前的日志文件执行归档
deleteBeforHour=72      # 多长时间以前的归档文件执行删除操作
params="{\"log_prefix\":\"${logPrefix}\",\"path\":\"${modulepwd}/logs\",\"rotate_path\":\"${modulepwd}/logs/rotate\",\"before_hour\":${rotateBeforHour},\"delete_hour\":${deleteBeforHour}}"

#尝试删除已经存在的
cronDelete="${modulepwd}/logrotate"
cronDelete2="${cronDelete//\//\\/}"
sed -i '/'''${cronDelete2}'''/d' /var/spool/cron/root

#重新加入
echo "add module logrotate"
echo "*/1 * * * * ${modulepwd}/logrotate -json '${params}' >/dev/null 2>&1" >>/var/spool/cron/root

