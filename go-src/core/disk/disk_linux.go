//go:build linux

package disk

import (
	"os/exec"
	"strconv"
	"strings"
)

func getDiskUsage(path string) (total, free uint64, err error) {
	out, err := exec.Command("df", "-B1", path).CombinedOutput()
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0, 0, nil
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0, 0, nil
	}
	total, _ = strconv.ParseUint(fields[1], 10, 64)
	free, _ = strconv.ParseUint(fields[3], 10, 64)
	return total, free, nil
}

func dirSize(path string) (uint64, error) {
	out, err := exec.Command("du", "-sb", path).CombinedOutput()
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return 0, nil
	}
	size, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
