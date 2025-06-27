package snowflake

import (
	"sync"
	"time"
)

var singletonMutex sync.Mutex
var idGenerator *defaultIdGenerator

// setIdGenerator 设置生成器
func setIdGenerator(options *idGeneratorOptions) {
	singletonMutex.Lock()
	idGenerator = newDefaultIdGenerator(options)
	singletonMutex.Unlock()
}

func extractTime(id int64) time.Time {
	return idGenerator.ExtractTime(id)
}

func NextId() int64 {
	if idGenerator == nil {
		panic("Please initialize Yitter.IdGeneratorOptions first.")
	}
	return idGenerator.NewLong()
}

func init() {
	var options = newIdGeneratorOptions(1)
	setIdGenerator(options)
}
