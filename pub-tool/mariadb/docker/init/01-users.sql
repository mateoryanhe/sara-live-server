-- 账号规范(与测试服一致)
-- root: 超级管理员, 仅 127.0.0.1 (+ localhost)
-- appuser: 超级管理员, 可远程

CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY 'Appledev882116';
CREATE USER IF NOT EXISTS 'appuser'@'%' IDENTIFIED BY 'Appledev882116';
CREATE USER IF NOT EXISTS 'appuser'@'localhost' IDENTIFIED BY 'Appledev882116';

ALTER USER 'root'@'localhost' IDENTIFIED BY 'Appledev882116';
ALTER USER 'root'@'127.0.0.1' IDENTIFIED BY 'Appledev882116';
ALTER USER 'appuser'@'%' IDENTIFIED BY 'Appledev882116';
ALTER USER 'appuser'@'localhost' IDENTIFIED BY 'Appledev882116';

GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'appuser'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'appuser'@'localhost' WITH GRANT OPTION;

CREATE DATABASE IF NOT EXISTS `live_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 删除远程 root(若存在)
DROP USER IF EXISTS 'root'@'%';

FLUSH PRIVILEGES;
