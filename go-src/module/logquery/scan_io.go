package logquery

import (
	"bufio"
	"os"
)

func openLogFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}

func newLogScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)
	return scanner
}

func scanLogFile(filePath string, fn func(line string) bool) error {
	file, err := openLogFile(filePath)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	defer file.Close()

	scanner := newLogScanner(file)
	for scanner.Scan() {
		if !fn(scanner.Text()) {
			break
		}
	}
	return scanner.Err()
}
