#!/bin/bash
# MariaDB 11.4 LTS 低内存安装 (1GB 机器, mysqld 目标约 80MB)
# Amazon Linux 2023 需指定 RHEL9 源: --os-type=rhel --os-version=9

set -euo pipefail

MARIADB_VERSION="${MARIADB_VERSION:-11.4}"
ROOT_PWD="${ROOT_PWD:-Appledev882116}"
APP_USER="${APP_USER:-appuser}"
APP_PWD="${APP_PWD:-Appledev882116}"
DB_NAME="${DB_NAME:-live_db}"

if [[ $EUID -ne 0 ]]; then
  exec sudo -E bash "$0" "$@"
fi

command -v curl >/dev/null 2>&1 || { echo "需要 curl(curl-minimal 即可)"; exit 1; }

if ! rpm -q MariaDB-server >/dev/null 2>&1; then
  curl -LsS https://r.mariadb.com/downloads/mariadb_repo_setup | \
    bash -s -- --mariadb-server-version="${MARIADB_VERSION}" --os-type=rhel --os-version=9 --skip-check-installed
  dnf clean all
  dnf makecache -y
  dnf install -y MariaDB-server MariaDB-client
fi

cat > /etc/my.cnf.d/99-lowmem.cnf <<EOF
[mysqld]
performance_schema = OFF
port = 13307
bind-address = 0.0.0.0
skip-name-resolve = 1
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init-connect = 'SET NAMES utf8mb4'
sql_mode = STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
innodb_buffer_pool_size = 32M
innodb_log_buffer_size = 512K
innodb_buffer_pool_instances = 1
innodb_flush_log_at_trx_commit = 2
innodb_file_per_table = 1
innodb_flush_method = O_DIRECT
innodb_strict_mode = 1
max_connections = 20
table_open_cache = 32
table_definition_cache = 200
thread_cache_size = 4
key_buffer_size = 8M
sort_buffer_size = 256K
read_buffer_size = 256K
read_rnd_buffer_size = 256K
join_buffer_size = 256K
tmp_table_size = 4M
max_heap_table_size = 4M
slow_query_log = 0
general_log = 0
[client]
default-character-set = utf8mb4
[mysql]
default-character-set = utf8mb4
EOF

systemctl enable mariadb
systemctl restart mariadb

for i in $(seq 1 30); do
  mariadb-admin ping --silent && break
  sleep 1
done

cat > /tmp/init_mariadb.sql <<SQL
-- root: 仅本机 127.0.0.1 (+ localhost socket 运维)
ALTER USER 'root'@'localhost' IDENTIFIED BY '${ROOT_PWD}';
CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY '${ROOT_PWD}';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;

-- appuser: 超级管理员, 可远程
CREATE USER IF NOT EXISTS '${APP_USER}'@'%' IDENTIFIED BY '${APP_PWD}';
CREATE USER IF NOT EXISTS '${APP_USER}'@'localhost' IDENTIFIED BY '${APP_PWD}';
GRANT ALL PRIVILEGES ON *.* TO '${APP_USER}'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO '${APP_USER}'@'localhost' WITH GRANT OPTION;

CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
FLUSH PRIVILEGES;
SQL

mariadb -uroot < /tmp/init_mariadb.sql 2>/dev/null || \
  mariadb -uroot --password="${ROOT_PWD}" < /tmp/init_mariadb.sql

rm -f /tmp/init_mariadb.sql

mariadb --version
mariadb -uroot --password="${ROOT_PWD}" -e \
  "SHOW VARIABLES WHERE Variable_name IN ('innodb_buffer_pool_size','max_connections','performance_schema','character_set_server');"

echo "安装完成"
echo "连接串: mysql:${APP_USER}:${APP_PWD}@tcp(127.0.0.1:13307)/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"
