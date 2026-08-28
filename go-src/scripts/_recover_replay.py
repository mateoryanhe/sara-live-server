"""Replay StrReplace/Write/Delete from agent transcript jsonl for given line ranges."""
import json
import sys
from pathlib import Path

ROOT = Path(r"D:/company-code/sara-live-server")
TRANSCRIPT = Path(
    r"C:/Users/hw/.cursor/projects/d-company-code-sara-live-server/agent-transcripts"
    r"/d05d4488-6aa4-423c-b540-696d875df8cc/d05d4488-6aa4-423c-b540-696d875df8cc.jsonl"
)


def rel_path(raw: str) -> Path:
    s = raw.replace("\\", "/")
    prefix = "d:/company-code/sara-live-server/"
    if s.lower().startswith(prefix):
        s = s[len(prefix):]
    return ROOT / s


def apply_range(start: int, end: int, label: str) -> tuple[int, int]:
    applied = 0
    skipped = 0
    print(f"\n=== {label} (lines {start}-{end}) ===")
    with TRANSCRIPT.open(encoding="utf-8") as f:
        for i, line in enumerate(f, 1):
            if i < start or i > end:
                continue
            try:
                obj = json.loads(line)
            except json.JSONDecodeError:
                continue
            for c in obj.get("message", {}).get("content", []):
                if c.get("type") != "tool_use":
                    continue
                name = c.get("name")
                data = c.get("input", {})
                if name == "Write":
                    path = rel_path(data["path"])
                    path.parent.mkdir(parents=True, exist_ok=True)
                    path.write_text(data["contents"], encoding="utf-8", newline="\n")
                    print(f"Write {path.relative_to(ROOT)}")
                    applied += 1
                elif name == "Delete":
                    path = rel_path(data["path"])
                    if path.exists():
                        path.unlink()
                        print(f"Delete {path.relative_to(ROOT)}")
                        applied += 1
                elif name == "StrReplace":
                    path = rel_path(data["path"])
                    if not path.exists():
                        print(f"SKIP missing file line {i}: {path.relative_to(ROOT)}", file=sys.stderr)
                        skipped += 1
                        continue
                    text = path.read_text(encoding="utf-8")
                    old, new = data["old_string"], data["new_string"]
                    if old not in text:
                        print(f"SKIP old-not-found line {i}: {path.relative_to(ROOT)}", file=sys.stderr)
                        skipped += 1
                        continue
                    count = text.count(old)
                    if count > 1 and not data.get("replace_all"):
                        print(f"SKIP multiple-{count} line {i}: {path.relative_to(ROOT)}", file=sys.stderr)
                        skipped += 1
                        continue
                    if data.get("replace_all"):
                        text = text.replace(old, new)
                    else:
                        text = text.replace(old, new, 1)
                    path.write_text(text, encoding="utf-8", newline="\n")
                    print(f"StrReplace {path.relative_to(ROOT)}")
                    applied += 1
    print(f"{label}: applied={applied} skipped={skipped}")
    return applied, skipped


def main():
    total_a = total_s = 0
    for start, end, label in [
        (6354, 6406, "9b80992"),
        (6420, 6727, "67e9cd2"),
    ]:
        a, s = apply_range(start, end, label)
        total_a += a
        total_s += s
    print(f"\nTOTAL: applied={total_a} skipped={total_s}")


if __name__ == "__main__":
    main()
