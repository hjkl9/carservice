package time

import "time"

func StdFormat1(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 自定义时间类
type WellTime struct {
	time.Time // 继承 time.Time
}

// StdFmt1 扩展时间格式化方法
func (wt WellTime) StdFmt1() string {
	return wt.Time.Format("2006-01-02 15:04:05")
}
