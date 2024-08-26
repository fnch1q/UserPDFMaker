package data

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func NewFile(path string) (*File, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file := &File{
		Path:       path,
		Name:       fileInfo.Name(),
		Size:       fileInfo.Size(),
		UpdateTime: fileInfo.ModTime(),
		Hash:       calculateCRC32(path),
	}

	return file, nil
}

func calculateCRC32(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	_, err = io.Copy(hash, file)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%08X", hash.Sum32())
}
