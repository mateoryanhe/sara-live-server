国旗 Icon 生成（pub-tool/country-flags）
========================================

作用
----
从 go-src/constants/country/countries.go 解析国家简码，
下载 PNG 国旗到本目录 flags/{code}.png。

来源
----
flagcdn.com（例: https://flagcdn.com/w80/id.png）
默认宽度 80，可在 generate.ps1 参数里改。

用法
----
1. 双击 一键生成.bat（或 generate.bat）生成 flags/
2. 将 flags 目录打成 zip
3. 在 CMS「配置 → 平台与资源 → 国旗资源部署」上传 zip 发布
   - 与头像同一套 upload 存储：本地 storagePath 或云桶(S3/R2)
   - 相对路径：country-flags/{version}/{code}.png
   - version 写入表 country_flag_cfgs
   - App 通过 GetUrlByName（与头像同一资源域名）访问
   - 自动清理旧 version（本地目录 + 云对象）

说明
----
- 本目录只负责本地生成 PNG，不上传服务器
- 改 countries.go 后重新生成，再打包上传 CMS 即可
