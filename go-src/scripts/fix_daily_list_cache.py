# -*- coding: utf-8 -*-
from pathlib import Path
import re

root = Path(__file__).resolve().parents[1] / "dao" / "statdao"
changed = []
for path in sorted(root.glob("daily_*.go")):
    text = path.read_text(encoding="utf-8")
    if "NewPermanentListCache" not in text or "func ListRecentDaily" not in text:
        continue
    orig = text

    text = re.sub(
        r"(var \(\n\t[^\n]+StatCacheMgr \*cache\.RowCache\[[^\]]+\])\n\t[^\n]+StatListCacheMgr \*cache\.ListCache\[[^\]]+\]\n\)",
        r"\1\n)",
        text,
        count=1,
    )
    text = re.sub(
        r"\n\t[a-zA-Z0-9_]+StatListCacheMgr = cache\.NewPermanentListCache\[[^\]]+\]\(\)",
        "",
        text,
        count=1,
    )
    text2, n = re.subn(
        r"\tif limit <= 0 \{\n\t\tlimit = 30\n\t\}\n"
        r'\tkey := fmt\.Sprintf\("%d:%s", limit, entity\.FormatDailyLoginStatDate\(time\.Now\(\)\)\)\n'
        r"\treturn \w+ListCacheMgr\.MustGetList\(gctx\.New\(\), key, func\(ctx context\.Context\) \(\[\]\*entity\.\w+, error\) \{\n"
        r"\t\treturn (loadRecentDaily\w+FromDB)\(limit\), nil\n"
        r"\t\}\)\n",
        "\tif limit <= 0 {\n\t\tlimit = 30\n\t}\n"
        "\t// CMS 按日图必须含「今天」:永久列表缓存在当天首条入库前命中后不会再刷新，会缺当天点。\n"
        "\treturn \\1(limit)\n",
        text,
        count=1,
    )
    if n != 1:
        print("SKIP list replace", path.name, n)
        continue
    text = text2
    if "fmt." not in text:
        text = text.replace('\t"fmt"\n', "")
    if "time." not in text:
        text = text.replace('\t"time"\n', "")
    text = re.sub(r"\n{3,}", "\n\n", text)
    if text == orig:
        print("SKIP no change", path.name)
        continue
    path.write_text(text, encoding="utf-8")
    changed.append(path.name)

print("changed", len(changed))
for name in changed:
    print(name)
