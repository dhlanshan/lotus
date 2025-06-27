package decorate

import "fmt"

// DataClasses 定义一个接口，包含要执行的方法
type DataClasses interface {
	postInit()
}

// NewStruct 包装一个带有初始化操作的struct
func NewStruct[T DataClasses](obj T) T {
	obj.postInit()
	fmt.Println("1")
	return obj
}
