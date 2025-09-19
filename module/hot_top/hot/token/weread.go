package token

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
)

// GetWereadID 获取微信读书的书籍ID
// 感谢 MCBBC 及 ChatGPT
func GetWereadID(bookId string) (string, error) {
	// 使用 MD5 哈希算法创建哈希对象
	hash := md5.New()
	hash.Write([]byte(bookId))
	str := hex.EncodeToString(hash.Sum(nil))

	// 取哈希结果的前三个字符作为初始值
	strSub := str[0:3]

	var fa []interface{}
	// 判断书籍 ID 的类型并进行转换
	if matched, _ := regexp.MatchString(`^\d*$`, bookId); matched {
		// 如果书籍 ID 只包含数字，则将其拆分成长度为 9 的子字符串，并转换为十六进制表示
		var chunks []string
		for i := 0; i < len(bookId); i += 9 {
			end := i + 9
			if end > len(bookId) {
				end = len(bookId)
			}
			chunk := bookId[i:end]
			num, err := strconv.ParseInt(chunk, 10, 64)
			if err != nil {
				return "", fmt.Errorf("解析数字块失败: %v", err)
			}
			chunks = append(chunks, fmt.Sprintf("%x", num))
		}
		fa = []interface{}{"3", chunks}
	} else {
		// 如果书籍 ID 包含其他字符，则将每个字符的 Unicode 编码转换为十六进制表示
		var hexStr string
		for _, r := range bookId {
			hexStr += fmt.Sprintf("%x", r)
		}
		fa = []interface{}{"4", []string{hexStr}}
	}

	// 将类型添加到初始值中
	strSub += fa[0].(string)

	// 将数字 2 和哈希结果的后两个字符添加到初始值中
	strSub += "2" + str[len(str)-2:]

	// 处理转换后的子字符串数组
	chunks := fa[1].([]string)
	for i, sub := range chunks {
		subLength := len(sub)
		subLengthHex := fmt.Sprintf("%x", subLength)
		// 如果长度只有一位数，则在前面添加 0
		if len(subLengthHex) == 1 {
			subLengthHex = "0" + subLengthHex
		}
		// 将长度和子字符串添加到初始值中
		strSub += subLengthHex + sub
		// 如果不是最后一个子字符串，则添加分隔符 'g'
		if i < len(chunks)-1 {
			strSub += "g"
		}
	}

	// 如果初始值长度不足 20，从哈希结果中取足够的字符补齐
	if len(strSub) < 20 {
		needed := 20 - len(strSub)
		strSub += str[0:needed]
	}

	// 使用 MD5 哈希算法创建新的哈希对象
	finalHash := md5.New()
	finalHash.Write([]byte(strSub))
	finalStr := hex.EncodeToString(finalHash.Sum(nil))

	// 取最终哈希结果的前三个字符并添加到初始值的末尾
	strSub += finalStr[0:3]

	return strSub, nil
}
