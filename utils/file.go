package utils

import (
	"fmt"
	"os"
	"path"
)

func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CopyWithOptions(srcDir, dstDir, prefixForDstFile string) error {
	srcFiles, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}
	if !Exists(dstDir) {
		err = os.MkdirAll(dstDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}
	for _, srcFile := range srcFiles {
		srcFilePath := path.Join(srcDir, srcFile.Name())
		dstFilePath := path.Join(dstDir, prefixForDstFile+srcFile.Name())
		if srcFile.IsDir() {
			if err := Copy(srcFilePath, dstFilePath); err != nil {
				return fmt.Errorf("failed to copy directory: %w", err)
			}
		} else {
			data, err := os.ReadFile(srcFilePath)
			if err != nil {
				return fmt.Errorf("failed to read source file: %w", err)
			}
			if err := os.WriteFile(dstFilePath, data, 0644); err != nil {
				return fmt.Errorf("failed to write destination file: %w", err)
			}
		}
	}
	return nil
}

func Copy(srcDir, dstDir string) error {
	return CopyWithOptions(srcDir, dstDir, "")
}

func ListFiles(input string, extension string, depth int) ([]string, error) {
	f, err := os.Stat(input)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		if path.Ext(input) == extension {
			return []string{input}, nil
		}
		return nil, nil
	}
	if depth == 0 {
		return nil, nil
	}
	files, err := os.ReadDir(input)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, file := range files {
		list, err := ListFiles(path.Join(input, file.Name()), extension, depth-1)
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}
		out = append(out, list...)
	}
	return out, nil
}
