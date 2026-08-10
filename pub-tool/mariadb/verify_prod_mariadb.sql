SELECT User, Host FROM mysql.global_priv WHERE User IN ('root','appuser','') ORDER BY User, Host;
SHOW VARIABLES WHERE Variable_name IN ('innodb_buffer_pool_size','innodb_log_file_size','max_connections','bind_address','port');
