package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PathExists 判断目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GenDir(wantDir string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取根目录失败")
	}

	staticPath := filepath.Join(dir, wantDir)
	if _, err = os.Stat(staticPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(staticPath, os.ModePerm); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return staticPath, nil
}

func FileMd5(f io.Reader) (md5Str string, err error) {
	hash := md5.New()
	if _, err = io.Copy(hash, f); err != nil {
		err = fmt.Errorf("计算文件hash失败")
		return
	}
	md5Str = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

func GetFileExtension(fName string) string {
	index := strings.LastIndex(fName, ".")
	if index == -1 {
		return ""
	}

	return fName[index+1:]
}
