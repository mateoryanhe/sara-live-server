#!/usr/bin/env bash
# 筛选 GoFrame access 日志中状态码非 200 的请求,按执行时间降序输出
#
# 日志格式示例:
#   2026-07-06T00:07:44.278 {bf60829...} 404 "GET https cms.bigtktool.shop /robots.txt HTTP/1.1" 0.000, 216.73.216.54, "", "Mozilla/5.0 ..."
#
# 用法:
#   ./filter_access_non200.sh [日志路径] [目标时间] [期望状态码]
#
# 参数:
#   日志路径     文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
#   目标时间     可选,匹配行首时间前缀,如 2026-07-06 / 2026-07-06T02
#   期望状态码   可选,默认 200;输出所有 status != 期望状态码 的记录
#
# 示例:
#   ./filter_access_non200.sh access-2026-07-06.log
#   ./filter_access_non200.sh 07-06
#   ./filter_access_non200.sh docs/2026-07
#   ./filter_access_non200.sh docs 2026-07-06T02
#   ./filter_access_non200.sh "access-*.log" 2026-07-06T02
#   zcat access-*.log.gz | ./filter_access_non200.sh - 2026-07-06

set -euo pipefail

LOG_PATH="${1:-access-*.log}"
TIME_FILTER="${2:-}"
EXPECT_STATUS="${3:-200}"

usage() {
	cat <<'EOF'
筛选 GoFrame access 日志中状态码非 200 的请求,按执行时间降序输出

用法:
  ./filter_access_non200.sh [日志路径] [目标时间] [期望状态码]

参数:
  日志路径     文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
  目标时间     可选,匹配行首时间前缀
  期望状态码   可选,默认 200

示例:
  ./filter_access_non200.sh access-2026-07-06.log
  ./filter_access_non200.sh 07-06
  ./filter_access_non200.sh docs/2026-07
  ./filter_access_non200.sh docs 2026-07-06T02
EOF
	exit 1
}

if [[ "$LOG_PATH" == "-h" || "$LOG_PATH" == "--help" ]]; then
	usage
fi

if ! [[ "$EXPECT_STATUS" =~ ^[0-9]+$ ]]; then
	echo "期望状态码必须是数字: $EXPECT_STATUS" >&2
	exit 1
fi

collect_log_files() {
	local path="$1"
	local -a files=()
	local search_dir="."
	local pattern=""

	if [[ "$path" == "-" ]]; then
		return 0
	fi

	if [[ -f "$path" ]]; then
		files=("$path")
	elif [[ -d "$path" ]]; then
		while IFS= read -r file; do
			files+=("$file")
		done < <(find "$path" -maxdepth 1 -name 'access-*.log' -type f | sort)
	else
		shopt -s nullglob
		local expanded=( $path )
		shopt -u nullglob

		if ((${#expanded[@]} > 0)) && { [[ -f "${expanded[0]}" ]] || [[ "$path" == *'*'* || "$path" == *'?'* ]]; }; then
			files=( "${expanded[@]}" )
		else
			if [[ "$path" == */* ]]; then
				search_dir="${path%/*}"
				pattern="${path##*/}"
				[[ ! -d "$search_dir" ]] && search_dir="."
			else
				pattern="$path"
				search_dir="."
			fi

			while IFS= read -r file; do
				files+=("$file")
			done < <(
				find "$search_dir" -maxdepth 1 -name 'access-*.log' -type f 2>/dev/null | sort | while IFS= read -r file; do
					name="${file##*/}"
					[[ "$name" == *"$pattern"* ]] && echo "$file"
				done
			)
		fi
	fi

	if ((${#files[@]} == 0)); then
		echo "未找到 access 日志文件: $path" >&2
		exit 1
	fi

	printf '%s\0' "${files[@]}"
}

filter_access_logs() {
	local expect_status="$1"
	local time_filter="$2"
	shift 2

	awk -v expect_status="$expect_status" -v time_filter="$time_filter" '
	function extract_path(request,    parts, n) {
		n = split(request, parts, " ")
		if (n >= 4) {
			return parts[4]
		}
		return request
	}

	function extract_method(request,    parts) {
		split(request, parts, " ")
		return parts[1]
	}

	{
		if (time_filter != "" && index($0, time_filter) != 1) {
			next
		}
		if (match($0, /^([0-9T:.-]+) \{([^}]+)\} ([0-9]+) "([^"]+)" ([0-9]+\.[0-9]+), ([^,]+), ([^,]*), "([^"]*)"/, m)) {
			status = m[3] + 0
			if (status == expect_status + 0) {
				next
			}
			duration_ms = m[5] * 1000
			path = extract_path(m[4])
			method = extract_method(m[4])
			printf "%.0f\t%s\t%s\t%s\t%s\t%.0fms\t%s\t%s\t%s\t%s\n", \
				duration_ms, m[1], m[2], status, method, duration_ms, path, m[6], m[8], m[4]
		}
	}
	' "$@"
}

echo "耗时,time,reqId,status,method,url,ip,userAgent"

run_filter() {
	filter_access_logs "$EXPECT_STATUS" "$TIME_FILTER" "$@" | sort -t $'\t' -k1,1nr | format_output
}

format_output() {
	while IFS=$'\t' read -r sort_ms req_time req_id status method duration path ip user_agent request; do
		[[ -z "$sort_ms" ]] && continue
		echo "${duration},${req_time},${req_id},${status},${method},${path},${ip},\"${user_agent}\""
	done
}

if [[ "$LOG_PATH" == "-" ]]; then
	run_filter
elif mapfile -d '' -t log_files < <(collect_log_files "$LOG_PATH"); then
	run_filter "${log_files[@]}"
fi
