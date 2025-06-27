package snowflake

import (
	"fmt"
	"strconv"
)

type snowWorkerM2 struct {
	*snowWorkerM1
}

func newSnowWorkerM2(options *idGeneratorOptions) iSnowWorker {
	return &snowWorkerM2{
		newSnowWorkerM1(options).(*snowWorkerM1),
	}
}

func (m2 snowWorkerM2) NextId() int64 {
	m2.Lock()
	defer m2.Unlock()
	currentTimeTick := m2.GetCurrentTimeTick()
	if m2._LastTimeTick == currentTimeTick {
		m2._CurrentSeqNumber++
		if m2._CurrentSeqNumber > m2.MaxSeqNumber {
			m2._CurrentSeqNumber = m2.MinSeqNumber
			currentTimeTick = m2.GetNextTimeTick()
		}
	} else {
		m2._CurrentSeqNumber = m2.MinSeqNumber
	}
	if currentTimeTick < m2._LastTimeTick {
		fmt.Println("Time error for {0} milliseconds", strconv.FormatInt(m2._LastTimeTick-currentTimeTick, 10))
	}
	m2._LastTimeTick = currentTimeTick
	result := int64(currentTimeTick<<m2._TimestampShift) + int64(m2.WorkerId<<m2.SeqBitLength) + int64(m2._CurrentSeqNumber)
	return result
}
