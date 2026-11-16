SELECT id, name, LEFT(pwd, 20) AS pwd_preview, status FROM cms_users ORDER BY id;

UPDATE cms_users SET pwd = '123456', updated_at = NOW(3);

SELECT id, name, pwd, status FROM cms_users ORDER BY id;