package disk

// GetDiskUsage 获取 path 所在分区的总容量与空闲容量(字节); 非 Linux 返回 0
func GetDiskUsage(path string) (total, free uint64, err error) {
	return getDiskUsage(path)
}

// DirSize 统计目录占用字节数; 非 Linux 返回 0
func DirSize(path string) (uint64, error) {
	return dirSize(path)
}
