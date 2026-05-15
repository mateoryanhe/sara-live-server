# MySQL一键安装脚本

此脚本用于在Linux服务器上一键安装MySQL 8.4.7，适用于6核16GB内存的生产服务器。

## 功能特性

- 从官方源下载MySQL 8.4.7二进制包：https://cdn.mysql.com//Downloads/MySQL-8.4/mysql-8.4.7-linux-glibc2.28-x86_64.tar.xz
- 自动创建数据目录、日志目录和临时目录
- 针对6核16GB内存服务器进行性能优化
- 创建安全的用户账户(root仅限本地访问，应用用户可远程访问)
- 自动配置防火墙规则
- 创建systemd服务文件以便管理

## 安装要求

- Linux操作系统
- 至少2GB可用内存
- 至少2GB可用磁盘空间
- root或sudo权限
- 网络连接(用于下载MySQL)

## 使用方法

```bash
# 以root或sudo用户身份运行安装脚本
sudo ./install_mysql.sh

# 或者指定自定义参数
sudo ./install_mysql.sh \
  --basedir /opt/mysql \
  --port 3307 \
  --root-password MySecurePassword123 \
  --app-user myapp \
  --app-password MyAppPassword123
```

### 参数说明

- `--basedir`: MySQL基础目录 (默认: /home/mysql)
- `--port`: MySQL端口 (默认: 11036)
- `--root-password`: root用户密码 (默认: Root123456!)
- `--app-user`: 应用用户名称 (默认: appuser)
- `--app-password`: 应用用户密码 (默认: AppUser123456!)

## 默认配置

- 数据目录: `/home/mysql/data`
- 日志目录: `/home/mysql/log`
- 套接字文件: `/home/mysql/mysql/mysql.sock`
- 配置文件: `/home/mysql/mysql/my.cnf`

## 用户权限

- root用户: 仅限本地访问(localhost, 127.0.0.1, ::1)，拥有所有权限
- 应用用户: 可从任何主机访问，拥有所有权限

## 故障排除

### MySQL服务无法启动

- 检查日志文件：`/home/mysql/log/mysqld.err`
- 确保端口11036未被占用
- 确保有足够的磁盘空间和内存

### 连接问题

- 如果出现 `ERROR 2002 (HY000): Can't connect to local MySQL server through socket` 错误：
    - 检查配置文件中的socket路径：`/home/mysql/mysql/mysql.sock`
    - 使用完整路径连接：`/home/mysql/mysql/bin/mysql -u root -p -S /home/mysql/mysql/mysql.sock`
    - 或者使用配置文件连接：`/home/mysql/mysql/bin/mysql --defaults-file=/home/mysql/mysql/my.cnf -u root -p`

### 用户权限问题

- root用户只能从本地连接（localhost, 127.0.0.1, ::1）
- 应用用户可以从任何地方连接
- 如需修改权限，请使用MySQL命令行工具

### 其他问题

- 如果安装过程中断，请删除 `/home/mysql` 目录并重新开始
- 检查防火墙设置是否允许端口11036
- 确保系统满足最低资源要求（至少2GB内存可用）
- 下载失败时检查网络连接，脚本会自动尝试使用curl或wget下载

## Logs Location

- Error log: `/home/mysql/log/mysqld.err`
- Slow query log: `/home/mysql/log/mysql-slow.log`
- General log: `/home/mysql/log/mysql.log`
- Initialization log: `/home/mysql/log/init.log`

/home/mysql/mysql/bin/mysql -h 127.0.0.1 -P 11036 -u root -p

nohup /home/mysql/mysql/bin/mysqld --defaults-file=/home/mysql/mysql/my.cnf --user=mysql --datadir=/home/mysql/data  >>
/dev/null 2>&1 &