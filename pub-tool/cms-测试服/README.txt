Hangzhou Server CMS Upload Instructions

One-click build and upload script:
   upload.bat - Double-click to execute, builds the Vue project and uploads the output to server
   The script will:
   1. Build the Vue project (D:\go-project\xrgameserver\cms) with dev environment
   2. Compress the built files to a zip file
   3. Upload it to hzaicoin.fun server and extract to /root/p-cms directory
   This approach reduces upload time and ensures latest build is deployed
   Note: plink.exe and pscp.exe (PuTTY tools) must be in the current directory

Configuration:
   All configuration parameters (server IP, username, password, directories, build settings, etc.) are defined directly in the upload.bat script
   To modify configuration:
   - SERVER_IP: Server IP address
   - SERVER_USER: Server username
   - SERVER_PASSWORD: Server password
   - SERVER_PORT: SSH port
   - VUE_PROJECT_DIR: Vue project directory to build
   - BUILD_OUTPUT_DIR: Directory containing build output (typically same as project dir)
   - BUILD_ENV: Build environment (dev/prod/staging)
   - REMOTE_DIR: Target directory on server

Usage:
   1. Ensure plink.exe and pscp.exe are in the current directory
   2. Make sure Node.js and npm are installed and available in PATH
   3. The script will try to use 7-Zip for compression if available, otherwise PowerShell compression
   4. Double-click the upload.bat file to execute the build and upload
   5. No manual password input required, the script will automatically use the preset password

Note: Ensure PuTTY tools (plink.exe and pscp.exe) are in the current directory before use