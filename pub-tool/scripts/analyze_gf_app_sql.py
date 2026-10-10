#!/usr/bin/env python3
"""Parse GoFrame database logger files under gf_app and summarize SQL frequency."""
import glob
import re
import sys
from collections import Counter, defaultdict

ENTRY_RE = re.compile(
    r"^\d{4}-\d{2}-\d{2}T[\d:.]+Z \[DEBU\] \{([0-9a-f]+)\} "
    r"\[\s*(\d+) ms\] \[default\] \[live_db\] \[rows:([^\]]*)\]"
)


def normalize_sql(sql: str) -> str:
    s = " ".join(sql.split())
    s = re.sub(r"'[^']*'", "'?'", s)
    s = re.sub(r"\b\d+\b", "?", s)
    return s[:600]


def main() -> int:
    paths = sorted(glob.glob("/home/ec2-user/log/gf_app/2026-10-*.log"))
    if len(sys.argv) > 1:
        paths = sorted(sys.argv[1:])

    patterns: Counter[str] = Counter()
    table_hits: Counter[str] = Counter()
    slow: list[tuple[int, str, str]] = []
    total = 0

    for path in paths:
        current = None
        with open(path, encoding="utf-8", errors="ignore") as f:
            for line in f:
                m = ENTRY_RE.match(line)
                if m:
                    if current and current.get("sql"):
                        pat = normalize_sql(current["sql"])
                        patterns[pat] += 1
                        total += 1
                        fm = re.search(r"\bFROM\s+`?([a-z_0-9]+)`?", current["sql"], re.I)
                        if fm:
                            table_hits[fm.group(1).lower()] += 1
                        if current["ms"] >= 50:
                            slow.append((current["ms"], pat[:140], current["trace"]))
                    current = {
                        "trace": m.group(1),
                        "ms": int(m.group(2)),
                        "sql": "",
                    }
                    continue
                if current is not None:
                    stripped = line.strip()
                    if stripped and not stripped.startswith("2026-"):
                        current["sql"] += (" " if current["sql"] else "") + stripped
            if current and current.get("sql"):
                pat = normalize_sql(current["sql"])
                patterns[pat] += 1
                total += 1

    print(f"LOG_FILES {len(paths)}")
    print(f"TOTAL_QUERIES {total}")
    print("\n=== TOP 30 SQL (normalized) ===")
    for pat, cnt in patterns.most_common(30):
        print(f"{cnt:6d}  {pat[:220]}")

    print("\n=== TOP tables (FROM) ===")
    for t, cnt in table_hits.most_common(25):
        print(f"{cnt:6d}  {t}")

    print("\n=== Slow >= 50ms (top 20) ===")
    for ms, pat, tr in sorted(slow, reverse=True)[:20]:
        print(f"{ms:4d}ms trace={tr} {pat}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
