#!/usr/bin/env bash
# 筛选 GoFrame access 日志中执行时间超过阈值的请求,按耗时降序输出
#
# 日志格式示例:
#   2026-07-06T00:09:14.302 {69846390...} 200 "POST https www.bigtktool.shop /liveRoom/reportLiveStartStatus HTTP/1.1" 0.141, 113.109.204.237, "", "Dart/3.11 (dart:io)"
#
# 用法:
#   ./filter_slow_access_requests.sh [日志路径] [阈值ms] [目标时间]
#
# 参数:
#   日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
#   阈值ms     只输出执行时间大于该值的请求,默认 100
#   目标时间   可选,匹配行首时间前缀,如 2026-07-06 / 2026-07-06T15
#
# 示例:
#   ./filter_slow_access_requests.sh access-2026-07-06.log
#   ./filter_slow_access_requests.sh 07-06                 # 模糊匹配 access-2026-07-06.log
#   ./filter_slow_access_requests.sh docs/2026-07 100      # 在 docs 目录模糊查找
#   ./filter_slow_access_requests.sh "access-*.log" 100
#   ./filter_slow_access_requests.sh docs 200 2026-07-06T04
#   zcat access-*.log.gz | ./filter_slow_access_requests.sh - 100

set -euo pipefail

LOG_PATH="${1:-access-*.log}"
MIN_MS="${2:-100}"
TIME_FILTER="${3:-}"

usage() {
	cat <<'EOF'
筛选 GoFrame access 日志中执行时间超过阈值的请求,按耗时降序输出

用法:
  ./filter_slow_access_requests.sh [日志路径] [阈值ms] [目标时间]

参数:
  日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
  阈值ms     只输出执行时间大于该值的请求,默认 100
  目标时间   可选,匹配行首时间前缀

示例:
  ./filter_slow_access_requests.sh access-2026-07-06.log
  ./filter_slow_access_requests.sh 07-06
  ./filter_slow_access_requests.sh docs/2026-07 100
  ./filter_slow_access_requests.sh docs 200 2026-07-06T04
EOF
	exit 1
}

if [[ "$LOG_PATH" == "-h" || "$LOG_PATH" == "--help" ]]; then
	usage
fi

if ! [[ "$MIN_MS" =~ ^[0-9]+$ ]]; then
	echo "阈值ms 必须是数字: $MIN_MS" >&2
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
	local min_ms="$1"
	local time_filter="$2"
	shift 2

	awk -v min_ms="$min_ms" -v time_filter="$time_filter" '
	function extract_path(request,    parts, n) {
		n = split(request, parts, " ")
		if (n >= 4) {
			return parts[4]
		}
		return request
	}

	{
		if (time_filter != "" && index($0, time_filter) != 1) {
			next
		}
		if (match($0, /^([0-9T:.-]+) \{([^}]+)\} ([0-9]+) "([^"]+)" ([0-9]+\.[0-9]+), ([^,]+),/, m)) {
			duration_ms = m[5] * 1000
			if (duration_ms > min_ms + 0) {
				path = extract_path(m[4])
				printf "%.0f\t%s\t%s\t%s\t%.0fms\t%s\t%s\n", duration_ms, m[1], m[2], m[3], duration_ms, path, m[6]
			}
		}
	}
	' "$@"
}

echo "耗时,time,reqId,status,url,ip"

run_filter() {
	filter_access_logs "$MIN_MS" "$TIME_FILTER" "$@" | sort -t $'\t' -k1,1nr | format_output
}

format_output() {
	while IFS=$'\t' read -r sort_ms req_time req_id status duration path ip; do
		[[ -z "$sort_ms" ]] && continue
		echo "${duration},${req_time},${req_id},${status},${path},${ip}"
	done
}

if [[ "$LOG_PATH" == "-" ]]; then
	run_filter
elif mapfile -d '' -t log_files < <(collect_log_files "$LOG_PATH"); then
	run_filter "${log_files[@]}"
fi
