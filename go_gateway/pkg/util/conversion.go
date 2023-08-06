package util

import "strconv"

// StrToUint32 字符串转uint32
func StrToUint32(str string) (uint32, error) {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

// StrToUint32Def0 字符串转uint32，失败返回0
func StrToUint32Def0(str string) uint32 {
	result, _ := StrToUint32(str)
	return result
}
