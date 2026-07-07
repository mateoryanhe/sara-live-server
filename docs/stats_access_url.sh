#!/usr/bin/env bash
# 统计 GoFrame access 日志中各 URL 访问次数,输出 Top N
#
# 日志格式示例:
#   2026-07-06T00:05:44.335 {traceId} 200 "POST https host /liveRoom/reportLiveStartStatus HTTP/1.1" 0.000, ip, "", "ua"
#
# 用法:
#   ./stats_access_url.sh [日志路径] [TopN]
#
# 参数:
#   日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
#   TopN       输出前 N 条,默认 50
#
# 示例:
#   ./stats_access_url.sh access-2026-07-06.log
#   ./stats_access_url.sh docs
#   ./stats_access_url.sh docs/07-06 20

set -euo pipefail

LOG_PATH="${1:-access-*.log}"
TOP_N="${2:-50}"

usage() {
	cat <<'EOF'
统计 GoFrame access 日志中各 URL 访问次数,按次数降序输出 Top N

用法:
  ./stats_access_url.sh [日志路径] [TopN]

参数:
  日志路径   文件 / 目录 / glob / 模糊关键字 / -,默认 access-*.log
  TopN       输出前 N 条,默认 50

示例:
  ./stats_access_url.sh access-2026-07-06.log
  ./stats_access_url.sh docs
  ./stats_access_url.sh docs/07-06 20
EOF
	exit 1
}

if [[ "$LOG_PATH" == "-h" || "$LOG_PATH" == "--help" ]]; then
	usage
fi

if ! [[ "$TOP_N" =~ ^[0-9]+$ ]] || (( TOP_N <= 0 )); then
	echo "TopN 必须是大于 0 的数字: $TOP_N" >&2
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

count_url_access() {
	awk '
	function parse_request(line,    request, parts, method, path) {
		if (!match(line, /"[A-Z]+ [^"]+"/)) {
			return ""
		}
		request = substr(line, RSTART + 1, RLENGTH - 2)
		split(request, parts, " ")
		if (length(parts) < 4) {
			return ""
		}
		method = parts[1]
		path = parts[4]
		if (path == "") {
			return ""
		}
		return method "\t" path
	}

	{
		key = parse_request($0)
		if (key != "") {
			count[key]++
			total++
		}
	}

	END {
		print "TOTAL\t" total
		for (key in count) {
			print count[key] "\t" key
		}
	}
	' "$@"
}

run_stats() {
	local stats_output total
	stats_output=$(count_url_access "$@")
	total=$(echo "$stats_output" | awk -F'\t' '$1=="TOTAL"{print $2; exit}')
	if [[ -z "$total" || "$total" == "0" ]]; then
		echo "未解析到有效 URL 访问记录" >&2
		exit 1
	fi

	echo "总请求数: ${total}, 独立URL数: $(echo "$stats_output" | awk -F'\t' '$1!="TOTAL"{c++} END{print c+0}')"
	echo "排名,method,url,访问次数,占比"
	echo "$stats_output" \
		| awk -F'\t' '$1!="TOTAL"{print}' \
		| sort -t $'\t' -k1,1nr \
		| head -n "$TOP_N" \
		| awk -v total="$total" -F'\t' '{
			rank = NR
			pct = (total > 0) ? ($1 * 100 / total) : 0
			printf "%d,%s,%s,%s,%.2f%%\n", rank, $2, $3, $1, pct
		}'
}

if [[ "$LOG_PATH" == "-" ]]; then
	run_stats
elif mapfile -d '' -t log_files < <(collect_log_files "$LOG_PATH"); then
	run_stats "${log_files[@]}"
fi
