package eosc

// IRequires 管理和操作依赖关系
type IRequires interface {
	// Set 设置一个ID及其依赖的ID列表
	Set(id string, requires []string)

	// Del 删除一个ID及其依赖关系
	Del(id string)

	// RequireByCount 获取依赖某个ID的ID的数量
	RequireByCount(requireId string) int

	// Requires 获取一个ID依赖的所有ID列表
	Requires(id string) []string

	// RequireBy 获取依赖某个ID的所有ID列表
	RequireBy(requireId string) []string
}
