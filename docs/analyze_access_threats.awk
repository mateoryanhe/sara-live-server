BEGIN {
	IGNORECASE = 1
}

function extract_host(request, parts) {
	split(request, parts, " ")
	return (length(parts) >= 3) ? parts[3] : ""
}

function extract_method(request, parts) {
	split(request, parts, " ")
	return parts[1]
}

function extract_path(request, parts) {
	split(request, parts, " ")
	return (length(parts) >= 4) ? parts[4] : request
}

function is_ip_host(host) {
	return host ~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/
}

function classify(path, host, method, status, ua, referer,    lower_path, lower_ua) {
	lower_path = tolower(path)
	lower_ua = tolower(ua)

	if (lower_path ~ /\/\.git|\/\.env|\/\.svn|\/\.hg|\/\.bzr|\/web\.config|\/\.aws/) {
		return "high|git-leak|源码/配置泄露探测"
	}
	if (lower_path ~ /\.php$|\/wp-admin|\/wp-login|\/phpmyadmin|\/phpunit|\/thinkphp|\/shell|\/cmd\.|\/eval-stdin|\/manager\/html|\/jmx-console|\/invoker/) {
		return "high|webshell-cms|Web/CMS漏洞探测"
	}
	if (lower_path ~ /\/owa\/|\/ecp\/|\/autodiscover|\/aspnet_client/) {
		return "high|exchange-scan|Exchange/OWA漏洞扫描"
	}
	if (lower_path ~ /\.\.\/|%2e%2e|%2f%2e%2e|\/etc\/passwd|\/proc\//) {
		return "high|path-traversal|路径穿越攻击"
	}
	if (lower_path ~ /union[ +]+select|or[ +]+1=1|<script|\/etc\/shadow/) {
		return "high|sql-injection|SQL注入/XSS探测"
	}
	if (lower_ua ~ /zgrab|scanner\/|nikto|masscan|nmap|sqlmap|censysinspect|shodan|acunetix|netsparker|wpscan|dirbuster|gobuster|nuclei/) {
		return "high|scanner-tool|自动化扫描工具"
	}
	if (lower_ua ~ /banner detection|ipip\.net/) {
		return "medium|banner-probe|服务指纹/Banner探测"
	}
	if (is_ip_host(host)) {
		return "medium|direct-ip|直接访问IP(非域名)"
	}
	if (path == "/" && (status == 403 || status == 404) && referer ~ /google\.com/) {
		return "medium|root-probe|根路径批量探测(可疑Referer)"
	}
	if (path == "/" && (status == 403 || status == 404) && lower_ua !~ /dart\//) {
		return "medium|root-probe|根路径探测"
	}
	if (lower_ua ~ /claudebot|claude-searchbot|gptbot|bytespider|petalbot/) {
		return "low|ai-crawler|AI爬虫"
	}
	if (lower_path ~ /\/robots\.txt$|\/sitemap\.xml$/) {
		return "low|seo-crawler|SEO/爬虫资源"
	}
	if (path ~ /^\/upload\/images$|^\/upload$/) {
		return "info|dir-listing|目录列表尝试(通常已被拒绝)"
	}
	if (path ~ /^\/cms/ || host ~ /cms\./) {
		return "info|cms-access|CMS后台/静态资源"
	}
	if (lower_ua ~ /dart\// || path ~ /^\/liveRoom\/|^\/userInfo\/|^\/gift\/|^\/shortVideo\//) {
		return "info|app-api|App正常API"
	}
	if (path ~ /^\/upload\/|^\/cms\/assets\//) {
		return "info|static-resource|静态资源"
	}
	return "info|normal|普通访问"
}

function parse_line(line) {
	if (match(line, /^([0-9T:.-]+) \{([^}]+)\} ([0-9]+) "([^"]+)" ([0-9]+\.[0-9]+), ([^,]+), ([^,]*), "([^"]*)"/, m)) {
		req_time = m[1]
		trace_id = m[2]
		status = m[3] + 0
		request = m[4]
		duration_ms = int(m[5] * 1000 + 0.5)
		ip = m[6]
		referer = m[7]
		ua = m[8]
		return 1
	}
	if (match(line, /^([0-9T:.-]+) ([0-9]+) "([^"]+)" ([0-9]+\.[0-9]+), ([^,]+), ([^,]*), "([^"]*)"/, n)) {
		req_time = n[1]
		trace_id = "-"
		status = n[2] + 0
		request = n[3]
		duration_ms = int(n[4] * 1000 + 0.5)
		ip = n[5]
		referer = n[6]
		ua = n[7]
		return 1
	}
	return 0
}

function threat_rank(level) {
	if (level == "high") return 1
	if (level == "medium") return 2
	if (level == "low") return 3
	return 4
}

function should_output(level) {
	if (mode == "all") return 1
	if (mode == "threat" || mode == "summary") {
		return (level == "high" || level == "medium" || level == "low")
	}
	return 0
}

{
	if (time_filter != "" && index($0, time_filter) != 1) {
		next
	}

	if (!parse_line($0)) {
		next
	}

	total++
	status_count[status]++

	host = extract_host(request)
	method = extract_method(request)
	path = extract_path(request)

	meta = classify(path, host, method, status, ua, referer)
	split(meta, parts, "|")
	level = parts[1]
	category = parts[2]
	desc = parts[3]

	category_count[category]++
	level_count[level]++

	if (!should_output(level)) {
		next
	}

	rank = threat_rank(level)
	printf "DETAIL\t%d\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%.0fms\t%s\n", \
		rank, level, category, req_time, trace_id, status, duration_ms, method, path, host, ip, ua, duration_ms, desc
}

END {
	print "STATS\tTOTAL\t" total
	for (s in status_count) {
		print "STATUS\t" s "\t" status_count[s]
	}
	for (l in level_count) {
		print "LEVEL\t" l "\t" level_count[l]
	}
	for (c in category_count) {
		print "CATEGORY\t" c "\t" category_count[c]
	}
}
