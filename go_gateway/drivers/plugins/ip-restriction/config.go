package ip_restriction

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrorConfigTypeError = errors.New("unknown ip list type")
	ErrorIpIllegal       = errors.New("[ip_restriction] Illegal IP")
)

// Config 插件配置
type Config struct {
	IPListType  string   `json:"ip_list_type" enum:"white,black" label:"列表类型"`
	IPWhiteList []string `json:"ip_white_list" label:"ip白名单列表"`
	IPBlackList []string `json:"ip_black_list" label:"ip黑名单列表"`
}

// doCheck 检查配置
func (c *Config) doCheck() error {
	if c.IPListType == "white" || c.IPListType == "black" {
		return nil
	}
	return ErrorConfigTypeError
}

// IPFilter 黑白名单业务逻辑
// @ip 客户端IP
// @return 是否继续往下执行 异常信息
type IPFilter func(ip string) (bool, error)

func (c *Config) genFilter() IPFilter {
	// 白名单
	if c.IPListType == "white" {
		return func(ip string) (bool, error) {
			flag := false // 默认非白名单，拒绝执行
			var err error // 异常
			for _, v := range c.IPWhiteList {
				// 放行所有
				if v == "*" {
					flag = true
					break
				}
				//
				v, err = convertIP(v)
				if err != nil {
					return false, err
				}
				// 客户端ip和配置的ip进行匹配
				if v != "" {
					if match(ip, v) {
						flag = true
						break
					}
				}
			}
			// 拒绝执行
			if !flag {
				return false, ErrorIpIllegal
			}
			// 命中白名单
			return true, nil
		}
	}
	// 黑名单
	if c.IPListType == "black" {
		return func(ip string) (bool, error) {
			var err error
			for _, v := range c.IPBlackList {
				if v == "*" {
					return false, ErrorIpIllegal
				}
				v, err = convertIP(v)
				if err != nil {
					return false, err
				}
				if v != "" {
					if match(ip, v) {
						return false, ErrorIpIllegal
					}
				}
			}
			return true, nil
		}
	}
	return nil
}

func convertIP(ip string) (string, error) {
	ipr := strings.Split(ip, "/")
	errInfo := "[ip_restriction] Illegal ip:" + ip
	if len(ipr) > 0 {
		ips := strings.Split(ipr[0], ".")
		ipLen := len(ips)
		if firstIndex := strings.Index(ipr[0], "*"); firstIndex > 0 {
			if lastIndex := strings.LastIndex(ipr[0], "*"); firstIndex == lastIndex && ips[ipLen-1] == "*" {
				v := ""
				for i := 0; i < 4; i++ {
					if i < ipLen-1 {
						v += ips[i] + "."
					} else {
						v += "0"
						if i != 3 {
							v += "."
						}
					}
				}
				v += "/" + strconv.Itoa((ipLen-1)*8)
				return v, nil
			} else {
				return "", errors.New(errInfo)
			}
		} else {
			if ipLen < 4 {
				return "", errors.New(errInfo)
			}
			return ip, nil
		}
	} else {
		return "", errors.New(errInfo)
	}
}

func ip2binary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, _ := strconv.ParseUint(s, 10, 8)

		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

func match(ip, iprange string) bool {
	ipb := ip2binary(ip)
	ipr := strings.Split(iprange, "/")
	if len(ipr) < 2 {
		return ip == ipr[0]
	}
	masklen, err := strconv.ParseUint(ipr[1], 10, 32)
	if err != nil {
		return false
	}
	iprb := ip2binary(ipr[0])
	return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
}
