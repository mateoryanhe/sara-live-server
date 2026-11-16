UPDATE domain_site_mappings SET
  domain = 'test-cms.zele.live',
  cert_file = '/home/ec2-user/ssl/cf-origin/cf-origin.pem',
  key_file = '/home/ec2-user/ssl/cf-origin/cf-origin.key'
WHERE site_key = 'cms';

SELECT domain, site_key, root, cert_file FROM domain_site_mappings WHERE site_key = 'cms';