#!/bin/bash
set -euo pipefail

ROOT_PWD='c63eac559a03dece518e3eb7a601b30e'
OLD_ROOT_PWD='Appledev882116'
CONF='/etc/my.cnf.d/99-sara-live.cnf'
BACKUP="/etc/my.cnf.d/99-lowmem.cnf.bak.$(date +%Y%m%d%H%M%S)"

if [[ $EUID -ne 0 ]]; then
  exec sudo -E bash "$0" "$@"
fi

mariadb_cmd() {
  if mariadb -uroot -p"${OLD_ROOT_PWD}" -e 'SELECT 1' >/dev/null 2>&1; then
    mariadb -uroot -p"${OLD_ROOT_PWD}" "$@"
  elif mariadb -uroot -p"${ROOT_PWD}" -e 'SELECT 1' >/dev/null 2>&1; then
    mariadb -uroot -p"${ROOT_PWD}" "$@"
  else
    mariadb -uroot "$@"
  fi
}

if [[ -f /etc/my.cnf.d/99-lowmem.cnf ]]; then
  cp /etc/my.cnf.d/99-lowmem.cnf "${BACKUP}"
fi

cat > "${CONF}" <<'EOF'
[mysqld]
performance_schema = OFF
port = 13307
bind-address = 127.0.0.1
skip-name-resolve = 1
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init-connect = 'SET NAMES utf8mb4'
sql_mode = STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION

# 约 2GB 内存预算(7.6G 机器, Go 预留 ~5G)
innodb_buffer_pool_size = 1280M
innodb_log_buffer_size = 16M
innodb_log_file_size = 128M
innodb_flush_log_at_trx_commit = 2
innodb_file_per_table = 1
innodb_flush_method = O_DIRECT
innodb_strict_mode = 1

max_connections = 80
table_open_cache = 512
table_definition_cache = 400
thread_cache_size = 16
open_files_limit = 32768

key_buffer_size = 16M
sort_buffer_size = 256K
read_buffer_size = 256K
read_rnd_buffer_size = 512K
join_buffer_size = 256K
tmp_table_size = 32M
max_heap_table_size = 32M

slow_query_log = 1
slow_query_log_file = /var/log/mariadb/slow.log
long_query_time = 2
log_slow_verbosity = query_plan
general_log = 0

[client]
port = 13307
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4
EOF

mkdir -p /var/log/mariadb
chown mysql:mysql /var/log/mariadb 2>/dev/null || chown mariadb:mariadb /var/log/mariadb 2>/dev/null || true

if [[ -f /etc/my.cnf.d/99-lowmem.cnf ]]; then
  mv /etc/my.cnf.d/99-lowmem.cnf /etc/my.cnf.d/99-lowmem.cnf.disabled
fi

systemctl restart mariadb

for i in $(seq 1 30); do
  mariadb-admin ping --silent && break
  sleep 1
done

mariadb_cmd <<SQL
ALTER USER 'root'@'localhost' IDENTIFIED BY '${ROOT_PWD}';
ALTER USER 'root'@'127.0.0.1' IDENTIFIED BY '${ROOT_PWD}';
DROP USER IF EXISTS 'root'@'%';
FLUSH PRIVILEGES;
SQL

while IFS=$'\t' read -r u h; do
  [[ -z "${u}" ]] || continue
  mariadb -uroot -p"${ROOT_PWD}" -e "DROP USER IF EXISTS \`${u}\`@\`${h}\`;"
done < <(mariadb -uroot -p"${ROOT_PWD}" -N -e "SELECT User, Host FROM mysql.global_priv WHERE User='';")

mariadb -uroot -p"${ROOT_PWD}" -e 'FLUSH PRIVILEGES;'

cat > /root/.my.cnf <<EOF
[client]
user=root
password=${ROOT_PWD}
port=13307
EOF
chmod 600 /root/.my.cnf

echo '=== MariaDB tuned ==='
mariadb -uroot -p"${ROOT_PWD}" -e "SHOW VARIABLES WHERE Variable_name IN ('innodb_buffer_pool_size','innodb_log_file_size','max_connections','bind_address','port');"
echo
mariadb -uroot -p"${ROOT_PWD}" -e "SELECT User, Host FROM mysql.global_priv WHERE User='root' ORDER BY Host;"
echo
free -h
