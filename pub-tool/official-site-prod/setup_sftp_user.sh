#!/bin/bash
set -euo pipefail

# 正式服官网 SFTP 专用账号(仅上传/删除 /home/ec2-user/cdn/official-site,不可执行命令)
SFTP_USER="official-site"
CHROOT="/var/sftp/${SFTP_USER}"
TARGET="/home/ec2-user/cdn/official-site"
PUBKEY_FILE="/tmp/official-site-sftp.pub"

if [[ $EUID -ne 0 ]]; then
  exec sudo -E bash "$0" "$@"
fi

if [[ ! -f "${PUBKEY_FILE}" ]]; then
  echo "missing ${PUBKEY_FILE}"
  exit 1
fi

if ! id "${SFTP_USER}" >/dev/null 2>&1; then
  useradd -m -d "/home/${SFTP_USER}" -s /sbin/nologin "${SFTP_USER}"
fi

mkdir -p "${CHROOT}/files"
chown root:root "${CHROOT}"
chmod 755 "${CHROOT}"
chown "${SFTP_USER}:${SFTP_USER}" "${CHROOT}/files"
chmod 755 "${CHROOT}/files"

if ! mountpoint -q "${CHROOT}/files"; then
  mount --bind "${TARGET}" "${CHROOT}/files"
fi
if ! grep -qF "${CHROOT}/files" /etc/fstab; then
  echo "${TARGET} ${CHROOT}/files none bind 0 0" >> /etc/fstab
fi

# bind 挂载后实际文件仍属 ec2-user，需改为 SFTP 用户可写（静态文件 644/755 仍可供 Go 读取）
chown -R "${SFTP_USER}:${SFTP_USER}" "${TARGET}"
find "${TARGET}" -type d -exec chmod 755 {} \;
find "${TARGET}" -type f -exec chmod 644 {} \;

install -d -m 700 -o "${SFTP_USER}" -g "${SFTP_USER}" "/home/${SFTP_USER}/.ssh"
install -m 600 -o "${SFTP_USER}" -g "${SFTP_USER}" "${PUBKEY_FILE}" "/home/${SFTP_USER}/.ssh/authorized_keys"

cat > "/etc/ssh/sshd_config.d/${SFTP_USER}.conf" <<EOF
Match User ${SFTP_USER}
    ChrootDirectory ${CHROOT}
    ForceCommand internal-sftp -d /files
    AllowTcpForwarding no
    X11Forwarding no
    PermitTTY no
EOF

sshd -t
systemctl reload sshd

echo "OK user=${SFTP_USER} chroot=${CHROOT} target=${TARGET}"
