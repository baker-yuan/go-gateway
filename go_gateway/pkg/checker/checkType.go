package checker

// CheckType Checker类型
type CheckType int

const (
	CheckTypeEqual    CheckType = iota // 全等匹配Checker类型
	CheckTypePrefix                    // 前缀匹配Checker类型
	CheckTypeSuffix                    // 后缀匹配Checker类型
	CheckTypeSub                       // 子串匹配Checker类型
	CheckTypeNotEqual                  // 非等匹配Checker类型
	CheckTypeNone                      // 空值匹配Checker类型
	CheckTypeExist                     // 存在匹配Checker类型
	CheckTypeNotExist                  // 不存在匹配Checker类型
	CheckTypeRegular                   // 区分大小写的正则匹配Checker类型
	CheckTypeRegularG                  // 不区分大小写的正则匹配Checker类型
	CheckTypeAll                       // 任意匹配Checker类型
	CheckMultiple                      // 复合匹配
)
