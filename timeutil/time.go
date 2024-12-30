package timeutil

// 时间相关处理

type SecondToDHMSResult struct {
	Day    int64 `json:"day"`    // 天
	Hour   int64 `json:"hour"`   // 小时
	Minute int64 `json:"minute"` // 分钟
	Second int64 `json:"second"` // 秒
}

// SecondToDHMS 将秒转为-天-小时-分钟-秒
func SecondToDHMS(seconds int64) *SecondToDHMSResult {
	const (
		minute = 60
		hour   = 60 * minute
		day    = 24 * hour
	)
	days := seconds / day
	seconds %= day
	hours := seconds / hour
	seconds %= hour
	minutes := seconds / minute
	seconds %= minute

	data := &SecondToDHMSResult{Day: days, Hour: hours, Minute: minutes, Second: seconds}

	return data
}
