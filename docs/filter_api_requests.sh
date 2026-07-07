#!/usr/bin/env bash
# 筛选应用日志中处理时间超过阈值的 API 响应,按耗时降序输出
#
# 日志格式示例:
#   ... log_util.go:121: time=2026-07-06 00:16:44.316,reqId=...,正常响应,处理时间=141ms,url=...,ip=...,authId=...,响应数据=...
#
# 用法:
#   ./filter_api_requests.sh [日志路径] [阈值ms] [目标时间]
#
# 参数:
#   日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 *.log
#   阈值ms     只输出处理时间大于该值的响应,默认 100
#   目标时间   可选,匹配 time= 前缀,如 2026-07-06 / "2026-07-06 15"
#
# 示例:
#   ./filter_api_requests.sh 2026-07-06.log
#   ./filter_api_requests.sh 07-06                 # 模糊匹配 2026-07-06.log
#   ./filter_api_requests.sh docs/2026-07 100    # 在 docs 目录模糊查找
#   ./filter_api_requests.sh 2026-07-06.log 200
#   zcat app-*.log.gz | ./filter_api_requests.sh - 100

set -euo pipefail

LOG_PATH="${1:-*.log}"
MIN_MS="${2:-100}"
TIME_FILTER="${3:-}"

usage() {
	cat <<'EOF'
筛选应用日志中处理时间超过阈值的 API 响应,按耗时降序输出

用法:
  ./filter_api_requests.sh [日志路径] [阈值ms] [目标时间]

参数:
  日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 *.log
  阈值ms     只输出处理时间大于该值的响应,默认 100
  目标时间   可选,匹配 time= 前缀

示例:
  ./filter_api_requests.sh 2026-07-06.log
  ./filter_api_requests.sh 07-06
  ./filter_api_requests.sh docs/2026-07 100
  ./filter_api_requests.sh 2026-07-06.log 200
  ./filter_api_requests.sh docs 100 2026-07-06T03
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
		done < <(find "$path" -maxdepth 1 -name '*.log' ! -name 'access-*.log' -type f | sort)
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
				find "$search_dir" -maxdepth 1 -name '*.log' ! -name 'access-*.log' -type f 2>/dev/null | sort | while IFS= read -r file; do
					name="${file##*/}"
					[[ "$name" == *"$pattern"* ]] && echo "$file"
				done
			)
		fi
	fi

	if ((${#files[@]} == 0)); then
		echo "未找到应用日志文件: $path" >&2
		exit 1
	fi

	printf '%s\0' "${files[@]}"
}

parse_slow_api_logs() {
	local time_filter="$1"
	local min_ms="$2"
	shift 2

	awk -v time_filter="$time_filter" -v min_ms="$min_ms" '
	function field_value(text, key,    n, parts, i, item) {
		n = split(text, parts, ",")
		for (i = 1; i <= n; i++) {
			item = parts[i]
			if (index(item, key "=") == 1) {
				return substr(item, length(key) + 2)
			}
		}
		return ""
	}

	function message_of(text,    n, parts, i) {
		n = split(text, parts, ",")
		for (i = 2; i <= n; i++) {
			if (parts[i] !~ /=/) {
				return parts[i]
			}
		}
		return ""
	}

	/log_util\.go:/ && /reqId=/ && /处理时间=/ {
		pos = index($0, " time=")
		if (pos == 0) {
			next
		}
		text = substr($0, pos + 1)

		if (!match(text, /处理时间=([0-9]+)ms/, dm)) {
			next
		}
		duration_ms = dm[1] + 0
		if (duration_ms <= min_ms) {
			next
		}

		req_time = field_value(text, "time")
		if (time_filter != "" && req_time !~ ("^" time_filter)) {
			next
		}

		req_id = field_value(text, "reqId")
		url = field_value(text, "url")
		auth_id = field_value(text, "authId")
		ip = field_value(text, "ip")
		message = message_of(text)

		printf "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", \
			duration_ms, req_time, req_id, url, auth_id, ip, message, duration_ms "ms"
	}
	' "$@"
}

echo "处理时间,time,reqId,url,authId,ip,message"

run_parse() {
	parse_slow_api_logs "$TIME_FILTER" "$MIN_MS" "$@" \
		| sort -t $'\t' -k1,1nr -k2,2r \
		| format_output
}

format_output() {
	while IFS=$'\t' read -r sort_ms req_time req_id url auth_id ip message duration; do
		[[ -z "$sort_ms" ]] && continue
		echo "${duration},${req_time},${req_id},${url},${auth_id},${ip},${message}"
	done
}

if [[ "$LOG_PATH" == "-" ]]; then
	run_parse
elif mapfile -d '' -t log_files < <(collect_log_files "$LOG_PATH"); then
	run_parse "${log_files[@]}"
fi
