 #!/bin/bash

# MySQL 8.4 installation script for production server
# This script installs MySQL 8.4 with secure configuration optimized for 6 core 16GB server running Go programs

set -e

# Default configuration
MYSQL_VERSION="8.4.7"
MYSQL_TARBALL="mysql-${MYSQL_VERSION}-linux-glibc2.28-x86_64.tar.xz"
MYSQL_URL="https://cdn.mysql.com//Downloads/MySQL-8.4/${MYSQL_TARBALL}"
MYSQL_USER="mysql"
MYSQL_GROUP="mysql"
DEFAULT_BASE_DIR="/home/mysql"
DEFAULT_PORT="11036"
DEFAULT_ROOT_PWD="Root123456!"
DEFAULT_APP_USER="appuser"
DEFAULT_APP_PWD="AppUser123456!"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --basedir)
            BASE_DIR="$2"
            shift 2
            ;;
        --port)
            MYSQL_PORT="$2"
            shift 2
            ;;
        --root-password)
            ROOT_PWD="$2"
            shift 2
            ;;
        --app-user)
            APP_USER="$2"
            shift 2
            ;;
        --app-password)
            APP_PWD="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --basedir DIR          MySQL base directory (default: /home/mysql)"
            echo "  --port PORT            MySQL port (default: 3306)"
            echo "  --root-password PWD    MySQL root password (default: Root123456!)"
            echo "  --app-user USER        Application user with remote access (default: appuser)"
            echo "  --app-password PWD     Application user password (default: AppUser123456!)"
            echo "  -h, --help             Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Set defaults if not provided
BASE_DIR=${BASE_DIR:-$DEFAULT_BASE_DIR}
MYSQL_PORT=${MYSQL_PORT:-$DEFAULT_PORT}
ROOT_PWD=${ROOT_PWD:-$DEFAULT_ROOT_PWD}
APP_USER=${APP_USER:-$DEFAULT_APP_USER}
APP_PWD=${APP_PWD:-$DEFAULT_APP_PWD}

MYSQL_HOME="${BASE_DIR}/mysql"
MYSQL_DATA="${BASE_DIR}/data"
MYSQL_LOG="${BASE_DIR}/log"
MYSQL_TMP="${BASE_DIR}/tmp"
PID_FILE="${MYSQL_HOME}/mysqld.pid"
SOCKET_FILE="${MYSQL_HOME}/mysql.sock"

# Determine if running with sudo or as root
if [[ $EUID -eq 0 ]]; then
    print_warning "Running as root user"
    USE_SUDO=false
    SUDO_CMD=""
else
    if command -v sudo &> /dev/null; then
        USE_SUDO=true
        SUDO_CMD="sudo"
    else
        print_error "sudo command not found but required for administrative tasks"
        exit 1
    fi
fi

# Check if MySQL is already installed
if [[ -d "$MYSQL_HOME" ]]; then
    print_error "MySQL is already installed at $MYSQL_HOME"
    print_error "Remove the directory or choose a different basedir"
    exit 1
fi







# Create mysql user and group if not exist
if ! getent group "$MYSQL_GROUP" > /dev/null 2>&1; then
    print_status "Creating mysql group..."
    ${SUDO_CMD} groupadd "$MYSQL_GROUP"
fi

if ! id "$MYSQL_USER" &>/dev/null; then
    print_status "Creating mysql user..."
    ${SUDO_CMD} useradd -r -g "$MYSQL_GROUP" -s /bin/false "$MYSQL_USER"
fi

# Create directories
print_status "Creating directories..."
${SUDO_CMD} mkdir -p "$BASE_DIR"
${SUDO_CMD} mkdir -p "$MYSQL_HOME"
${SUDO_CMD} mkdir -p "$MYSQL_DATA"
${SUDO_CMD} mkdir -p "$MYSQL_LOG"
${SUDO_CMD} mkdir -p "$MYSQL_TMP"
# Set permissions after extracting MySQL to ensure mysql user can access everything
${SUDO_CMD} chown -R "$MYSQL_USER":"$MYSQL_GROUP" "$BASE_DIR"
# Ensure the mysql user has write access to the data directory
${SUDO_CMD} chmod 750 "$MYSQL_DATA"
${SUDO_CMD} chmod 750 "$MYSQL_LOG"
${SUDO_CMD} chmod 750 "$MYSQL_TMP"

# Download MySQL
print_status "Downloading MySQL $MYSQL_VERSION..."
cd /tmp
if [[ ! -f "$MYSQL_TARBALL" ]]; then
    if command -v curl &> /dev/null; then
        print_status "Using curl to download MySQL..."
        curl -L "$MYSQL_URL" -o "$MYSQL_TARBALL"
    elif command -v wget &> /dev/null; then
        print_status "Using wget to download MySQL..."
        wget "$MYSQL_URL" -O "$MYSQL_TARBALL"
    else
        print_error "Neither curl nor wget is available for downloading"
        exit 1
    fi
    
    # Verify that the download was successful
    if [[ ! -f "$MYSQL_TARBALL" ]]; then
        print_error "Download failed. Could not download MySQL from $MYSQL_URL"
        exit 1
    fi
    
    # Check if the downloaded file is a valid archive
    if [[ $(file -b --mime-type "$MYSQL_TARBALL") != "application/x-xz" ]] && [[ $(file -b --mime-type "$MYSQL_TARBALL") != "application/gzip" ]]; then
        print_error "Downloaded file is not a valid archive"
        exit 1
    fi
    
    print_status "MySQL $MYSQL_VERSION downloaded successfully"
else
    print_status "MySQL tarball already exists, skipping download"
fi

# Extract MySQL
print_status "Extracting MySQL..."
tar -xf "$MYSQL_TARBALL" --strip-components=1 -C "$MYSQL_HOME"

# Change ownership
${SUDO_CMD} chown -R "$MYSQL_USER":"$MYSQL_GROUP" "$MYSQL_HOME"

# Create MySQL configuration file
print_status "Creating MySQL configuration file..."

cat > /tmp/my.cnf << EOF
[mysqld]
# Basic settings
port = $MYSQL_PORT
datadir = $MYSQL_DATA
socket = $SOCKET_FILE
pid-file = $PID_FILE
user = $MYSQL_USER

# Logging
log-error = $MYSQL_LOG/mysqld.err
slow-query-log = 1
slow-query-log-file = $MYSQL_LOG/mysql-slow.log
long_query_time = 2
general_log = 0
general_log_file = $MYSQL_LOG/mysql.log

# Connection settings
max_connections = 50
max_connect_errors = 10
max_allowed_packet = 64M
table_open_cache = 400
external-locking = FALSE
max_heap_table_size = 6M

# Buffer settings for 16GB RAM
key_buffer_size = 2M
sort_buffer_size = 2M
read_buffer_size = 2M
read_rnd_buffer_size = 8M
myisam_sort_buffer_size = 3M
thread_cache_size = 8

# InnoDB settings optimized for 16GB RAM and 6 cores
innodb_buffer_pool_size = 10M
innodb_buffer_pool_instances = 2
innodb_data_file_path = ibdata1:10M:autoextend
innodb_flush_log_at_trx_commit = 2
innodb_log_buffer_size = 6M
# Using innodb_redo_log_capacity instead of deprecated innodb_log_file_size and innodb_log_files_in_group
innodb_redo_log_capacity = 1G
innodb_max_dirty_pages_pct = 9
innodb_lock_wait_timeout = 120
innodb_flush_method = O_DIRECT
innodb_io_capacity = 20
innodb_io_capacity_max = 40
innodb_read_io_threads = 2
innodb_write_io_threads = 2

# Performance settings
tmp_table_size = 1M
max_heap_table_size = 1M
bulk_insert_buffer_size = 1M
join_buffer_size = 1M
performance_schema = OFF
thread_stack = 192K
server-id = 1
log-bin = mysql-bin
# Using binlog_expire_logs_seconds instead of deprecated expire_logs_days
binlog_expire_logs_seconds = 172800  # 2 days in seconds
max_binlog_size = 1G  # Maximum size of a single binary log file
binlog_cache_size = 1M  # Size of the cache for binary log

# Security
local-infile = 0
mysqlx = 0  # Disable X Plugin to avoid binding to port 33060

# Character set
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init-connect = 'SET NAMES utf8mb4'


# Security settings for production
local-infile = 0

[client]
port = $MYSQL_PORT
socket = $SOCKET_FILE
default-character-set = utf8mb4

[mysqldump]
quick
max_allowed_packet = 16M
default-character-set = utf8mb4

[mysql]
no-auto-rehash
default-character-set = utf8mb4
EOF

${SUDO_CMD} mv /tmp/my.cnf "$MYSQL_HOME/my.cnf"
${SUDO_CMD} chown "$MYSQL_USER":"$MYSQL_GROUP" "$MYSQL_HOME/my.cnf"

# Ensure data directory is clean before initialization
print_status "Preparing data directory..."
${SUDO_CMD} rm -rf "$MYSQL_DATA"/*

# Initialize MySQL data directory
print_status "Initializing MySQL data directory..."
if [[ $USE_SUDO == true ]]; then
    ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysqld" --defaults-file="$MYSQL_HOME/my.cnf" --initialize-insecure --user="$MYSQL_USER" 2>&1 | tee -a "$MYSQL_LOG/init.log"
else
    "$MYSQL_HOME/bin/mysqld" --defaults-file="$MYSQL_HOME/my.cnf" --initialize-insecure --user="$MYSQL_USER" 2>&1 | tee -a "$MYSQL_LOG/init.log"
fi

# Wait for initialization to complete
print_status "Waiting for initialization to complete..."
sleep 10

# Ensure proper ownership after initialization
${SUDO_CMD} chown -R "$MYSQL_USER":"$MYSQL_GROUP" "$MYSQL_DATA"

# For MySQL 8.4+, we should use --initialize-insecure which creates the system tables automatically
# So we don't need to run mysql_upgrade separately
print_status "MySQL data directory initialized with system tables"

# Start MySQL server
print_status "Starting MySQL server..."
if [[ $USE_SUDO == true ]]; then
    # Using mysqld directly instead of mysqld_safe for better control in automated script
    ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysqld" --defaults-file="$MYSQL_HOME/my.cnf" --user="$MYSQL_USER" --datadir="$MYSQL_DATA" >> "$MYSQL_LOG/mysqld.log" 2>&1 &
else
    "$MYSQL_HOME/bin/mysqld" --defaults-file="$MYSQL_HOME/my.cnf" --user="$MYSQL_USER" --datadir="$MYSQL_DATA" >> "$MYSQL_LOG/mysqld.log" 2>&1 &
fi
MYSQLD_PID=$!

# Wait for MySQL to start
print_status "Waiting for MySQL to start..."
timeout=60
while [[ $timeout -gt 0 ]]; do
    # Check if MySQL process is running and listening on the specified port
    if ps -p $! > /dev/null 2>&1; then
        # Check if the port is being listened on
        if nc -z localhost $MYSQL_PORT 2>/dev/null || (command -v ss >/dev/null && ss -tuln | grep -q ":$MYSQL_PORT ") || (command -v netstat >/dev/null && netstat -tuln | grep -q ":$MYSQL_PORT "); then
            print_status "MySQL process is running and listening on port $MYSQL_PORT"
            break
        fi
    fi
    sleep 2
    ((timeout -= 2))
done

if [[ $timeout -le 0 ]]; then
    print_error "MySQL server failed to start"
    print_error "Check log file: $MYSQL_LOG/mysqld.log"
    print_error "Check error log: $MYSQL_LOG/mysqld.err"
    # Check if there are any mysqld processes running
    if pgrep mysqld > /dev/null; then
        print_warning "Found mysqld processes still running. Killing them..."
        pkill -f mysqld
        sleep 3
    fi
    exit 1
fi

# Give MySQL a bit more time to be ready for connections
sleep 5

# Wait for MySQL to be ready
print_status "Waiting for MySQL to be ready..."
timeout=60
while [[ $timeout -gt 0 ]]; do
    if [[ $USE_SUDO == true ]]; then
        if ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysqladmin" --defaults-file="$MYSQL_HOME/my.cnf" ping > /dev/null 2>&1; then
            break
        fi
    else
        if "$MYSQL_HOME/bin/mysqladmin" --defaults-file="$MYSQL_HOME/my.cnf" ping > /dev/null 2>&1; then
            break
        fi
    fi
    sleep 1
    ((timeout--))
done

if [[ $timeout -le 0 ]]; then
    print_error "MySQL server failed to start"
    print_error "Check log file: $MYSQL_LOG/mysqld.err"
    print_error "Also check general log: $MYSQL_LOG/mysqld.log"
    # Check if there are any mysqld processes running
    if pgrep mysqld > /dev/null; then
        print_warning "Found mysqld processes still running. Killing them..."
        pkill -f mysqld
        sleep 3
    fi
    exit 1
fi

print_status "MySQL server started successfully"

# Set root password and configure users
print_status "Configuring MySQL users and security settings..."

# Create temporary SQL file to configure MySQL
cat > /tmp/mysql_config.sql << EOF
-- Create application user with remote access
CREATE USER IF NOT EXISTS '$APP_USER'@'%' IDENTIFIED BY '$APP_PWD';
GRANT ALL PRIVILEGES ON *.* TO '$APP_USER'@'%';

-- Ensure root user can only login locally (in case it's not already configured)
CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY '$ROOT_PWD';
CREATE USER IF NOT EXISTS 'root'@'::1' IDENTIFIED BY '$ROOT_PWD';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'::1' WITH GRANT OPTION;

-- Remove anonymous users
DELETE FROM mysql.user WHERE User='';

-- Remove remote root access
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');

FLUSH PRIVILEGES;

-- Create test database
CREATE DATABASE IF NOT EXISTS test_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

# Execute the configuration SQL
print_status "Configuring MySQL users..."
# For MySQL 8.4+, we need to handle the authentication differently
# First, try to connect and set the root password
if [[ $USE_SUDO == true ]]; then
    # Try connecting with the default authentication
    ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" --connect-expired-password --execute="ALTER USER USER() IDENTIFIED BY '$ROOT_PWD';" || {
        # If that doesn't work, wait a bit and try again
        sleep 5
        ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" --connect-expired-password --execute="ALTER USER USER() IDENTIFIED BY '$ROOT_PWD';" || {
            print_error "Failed to set root password with --connect-expired-password"
            # Try using mysql_secure_installation approach
            print_status "Trying alternative method to set root password..."
            # Create a temporary file with the password change command
            echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '$ROOT_PWD';" > /tmp/change_root_pwd.sql
            ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" < /tmp/change_root_pwd.sql || {
                print_error "Failed to set root password with alternative method"
                exit 1
            }
            rm -f /tmp/change_root_pwd.sql
        }
    }
else
    "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" --connect-expired-password --execute="ALTER USER USER() IDENTIFIED BY '$ROOT_PWD';" || {
        sleep 5
        "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" --connect-expired-password --execute="ALTER USER USER() IDENTIFIED BY '$ROOT_PWD';" || {
            print_error "Failed to set root password with --connect-expired-password"
            # Try using mysql_secure_installation approach
            print_status "Trying alternative method to set root password..."
            echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '$ROOT_PWD';" > /tmp/change_root_pwd.sql
            "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root --password="" < /tmp/change_root_pwd.sql || {
                print_error "Failed to set root password with alternative method"
                exit 1
            }
            rm -f /tmp/change_root_pwd.sql
        }
    }
fi

# Now run the remaining configuration with the new password
if [[ $USE_SUDO == true ]]; then
    ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root -p"$ROOT_PWD" < /tmp/mysql_config.sql || {
        print_error "Failed to configure MySQL users"
        exit 1
    }
else
    "$MYSQL_HOME/bin/mysql" --defaults-file="$MYSQL_HOME/my.cnf" -u root -p"$ROOT_PWD" < /tmp/mysql_config.sql || {
        print_error "Failed to configure MySQL users"
        exit 1
    }
fi

# Clean up temporary files
rm -f /tmp/mysql_config.sql





print_status "MySQL installation completed successfully!"
print_status "MySQL is running on port $MYSQL_PORT"
print_status "Data directory: $MYSQL_DATA"
print_status "Log directory: $MYSQL_LOG"
print_status "Configuration file: $MYSQL_HOME/my.cnf"

print_status "MySQL users created:"
print_status "  - root (localhost only) with password: $ROOT_PWD"
print_status "  - $APP_USER (remote access enabled) with password: $APP_PWD"

print_status "To connect locally: $MYSQL_HOME/bin/mysql -u root -p'$ROOT_PWD'"
print_status "To connect remotely: mysql -h <server_ip> -u $APP_USER -p'$APP_PWD' -P $MYSQL_PORT"

# Verify installation
print_status "Verifying installation..."
if [[ $USE_SUDO == true ]]; then
    if ${SUDO_CMD} -u "$MYSQL_USER" "$MYSQL_HOME/bin/mysqladmin" --defaults-file="$MYSQL_HOME/my.cnf" -u root -p"$ROOT_PWD" status > /dev/null 2>&1; then
        print_status "MySQL is running and accessible"
    else
        print_error "Could not verify MySQL connection"
    fi
else
    if "$MYSQL_HOME/bin/mysqladmin" --defaults-file="$MYSQL_HOME/my.cnf" -u root -p"$ROOT_PWD" status > /dev/null 2>&1; then
        print_status "MySQL is running and accessible"
    else
        print_error "Could not verify MySQL connection"
    fi
fi

print_status "Installation completed!"

# Additional information about socket configuration
print_status "MySQL Socket Configuration:"
print_status "  - Socket file: $SOCKET_FILE"
print_status "  - To connect with socket: $MYSQL_HOME/bin/mysql -u root -p -S $SOCKET_FILE"
print_status "  - To connect with config: $MYSQL_HOME/bin/mysql --defaults-file=$MYSQL_HOME/my.cnf -u root -p"
print_status "  - To connect remotely: mysql -h <server_ip> -P $MYSQL_PORT -u $APP_USER -p"
 #!/bin/bash
