package service

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fastduck/treasure-doc/module/reverse_index/comm/gid"
	"fastduck/treasure-doc/module/reverse_index/index"
	"fmt"
	"hash"
)

func Index(content ...string) error {
	// 校验content是否为空
	if len(content) == 0 {
		return nil
	}
	for _, v := range content {
		words := index.GetSeg().Cut(v, true)
		ids := gid.BatchGenId(len(words))
		for i, w := range words {
			index.GetIndexCache().Set(w, ids[i])
			index.GetContentCache().Set(ids[i], v)
		}
	}
	return nil
}

func Search(keyword string) ([]string, error) {
	words := index.GetSeg().Cut(keyword, true)
	var results []string
	ids := index.GetIndexCache().Search(words...)
	results = index.GetContentCache().Get(ids...)
	return results, nil
}

func List() any {
	return index.GetIndexCache().List()
}

func ComputeHash(algorithm string, data []byte) (string, error) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}

	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
