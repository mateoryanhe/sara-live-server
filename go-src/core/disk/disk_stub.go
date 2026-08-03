//go:build !linux

package disk

func getDiskUsage(path string) (total, free uint64, err error) {
	return 0, 0, nil
}

func dirSize(path string) (uint64, error) {
	return 0, nil
}
