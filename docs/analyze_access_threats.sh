#!/usr/bin/env bash
# 分析 GoFrame access 日志,自动识别扫描/攻击/可疑请求
#
# 日志格式:
#   2026-07-06T00:07:44.278 {traceId} 200 "POST https www.bigtktool.shop /path HTTP/1.1" 0.000, ip, referer, "ua"
#
# 用法:
#   ./analyze_access_threats.sh [日志路径] [模式|目标时间] [目标时间|模式]
#
# 模式:
#   threat   仅输出威胁/可疑请求(默认)
#   summary  输出统计汇总 + 威胁明细
#   all      输出全部分类结果
#
# 示例:
#   ./analyze_access_threats.sh access-2026-07-06.log
#   ./analyze_access_threats.sh 07-06 summary
#   ./analyze_access_threats.sh docs/2026-07 summary
#   ./analyze_access_threats.sh docs summary
#   ./analyze_access_threats.sh docs 2026-07-06T02 summary

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
AWK_SCRIPT="${SCRIPT_DIR}/analyze_access_threats.awk"

LOG_PATH="${1:-access-*.log}"
ARG2="${2:-}"
ARG3="${3:-}"
TIME_FILTER=""
MODE="threat"

if [[ "$ARG2" == "threat" || "$ARG2" == "summary" || "$ARG2" == "all" ]]; then
	MODE="$ARG2"
	TIME_FILTER="$ARG3"
elif [[ "$ARG3" == "threat" || "$ARG3" == "summary" || "$ARG3" == "all" ]]; then
	TIME_FILTER="$ARG2"
	MODE="$ARG3"
else
	TIME_FILTER="$ARG2"
	if [[ -n "$ARG3" ]]; then
		MODE="$ARG3"
	fi
fi

usage() {
	cat <<'EOF'
分析 GoFrame access 日志,自动识别扫描/攻击/可疑请求

用法:
  ./analyze_access_threats.sh [日志路径] [模式|目标时间] [目标时间|模式]

日志路径:
  文件 / 目录 / glob / 模糊关键字 / -

模式:
  threat   仅输出威胁/可疑请求(默认)
  summary  输出统计汇总 + 威胁明细
  all      输出全部分类结果

示例:
  ./analyze_access_threats.sh access-2026-07-06.log
  ./analyze_access_threats.sh 07-06 summary
  ./analyze_access_threats.sh docs/2026-07 summary
  ./analyze_access_threats.sh docs summary
  ./analyze_access_threats.sh docs 2026-07-06T02 summary
EOF
	exit 1
}

if [[ "$LOG_PATH" == "-h" || "$LOG_PATH" == "--help" ]]; then
	usage
fi

if [[ "$MODE" != "threat" && "$MODE" != "summary" && "$MODE" != "all" ]]; then
	echo "模式必须是 threat / summary / all: $MODE" >&2
	exit 1
fi

if [[ ! -f "$AWK_SCRIPT" ]]; then
	echo "缺少 awk 脚本: $AWK_SCRIPT" >&2
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

analyze_logs() {
	local mode="$1"
	local time_filter="$2"
	shift 2
	awk -v mode="$mode" -v time_filter="$time_filter" -f "$AWK_SCRIPT" "$@"
}

print_summary_header() {
	local total="$1"
	shift
	local -a stats_lines=("$@")

	echo "========================================"
	echo "       Access 日志威胁分析报告"
	echo "========================================"
	echo "总请求数: ${total}"
	echo
	echo "--- 状态码分布 ---"
	for line in "${stats_lines[@]}"; do
		[[ "$line" == STATUS* ]] || continue
		echo "  ${line#STATUS	}" | tr '\t' ': '
	done
	echo
}

print_level_summary() {
	local -a lines=("$@")
	echo "--- 威胁等级分布 ---"
	for line in "${lines[@]}"; do
		[[ "$line" == LEVEL* ]] || continue
		local level count
		level=$(echo "$line" | cut -f2)
		count=$(echo "$line" | cut -f3)
		case "$level" in
			high)   echo "  [高危]   ${count}" ;;
			medium) echo "  [中危]   ${count}" ;;
			low)    echo "  [低危]   ${count}" ;;
			info)   echo "  [信息]   ${count}" ;;
		esac
	done
	echo
}

print_category_summary() {
	local -a lines=("$@")
	echo "--- 分类统计(全量) ---"
	for line in "${lines[@]}"; do
		[[ "$line" == CATEGORY* ]] || continue
		local category count
		category=$(echo "$line" | cut -f2)
		count=$(echo "$line" | cut -f3)
		printf "  %-18s %s\n" "${category}:" "${count}"
	done
	echo
}

format_detail() {
	while IFS=$'\t' read -r tag rank level category req_time trace_id status duration_ms method path host ip ua duration desc; do
		[[ "$tag" != "DETAIL" || -z "$rank" ]] && continue
		echo "${level},${category},${req_time},${status},${method},${path},${host},${ip},\"${ua}\",${duration},${desc}"
	done
}

run_analyze() {
	local raw_output detail_output
	local total=0
	local -a stats_lines=()
	local -a detail_lines=()

	raw_output=$(analyze_logs "$MODE" "$TIME_FILTER" "$@")

	while IFS= read -r line || [[ -n "$line" ]]; do
		case "$line" in
			STATS*)
				stats_lines+=("$line")
				total=$(echo "$line" | cut -f3)
				;;
			STATUS*|LEVEL*|CATEGORY*)
				stats_lines+=("$line")
				;;
			DETAIL*)
				detail_lines+=("$line")
				;;
		esac
	done <<EOF
${raw_output}
EOF

	if ((${#detail_lines[@]} > 0)); then
		detail_output=$(printf '%s\n' "${detail_lines[@]}" | sort -t $'\t' -k2,2n -k5,5r)
	else
		detail_output=""
	fi

	if [[ "$MODE" == "summary" ]]; then
		print_summary_header "$total" "${stats_lines[@]}"
		print_level_summary "${stats_lines[@]}"
		print_category_summary "${stats_lines[@]}"
		echo "--- 威胁/可疑明细(按等级+时间降序) ---"
	fi

	echo "等级,分类,time,status,method,url,host,ip,userAgent,耗时,说明"
	if [[ -n "$detail_output" ]]; then
		echo "$detail_output" | format_detail
	fi
}

if [[ "$LOG_PATH" == "-" ]]; then
	run_analyze
elif mapfile -d '' -t log_files < <(collect_log_files "$LOG_PATH"); then
	run_analyze "${log_files[@]}"
fi
