package time

import "time"

const (
	FormatYmd    = "2006-01-02"
	FormatYmdH   = "2006010215"
	FormatYmdHms = "2006-01-02 15:04:05"
	FormatYmdHm  = "2006-01-02 15:04"
	FormatUnix   = "2006-01-02T15:04:05"

	LayOutUnixNamo   string = "2006-01-02T15:04:05.000000+00:00"
	LayOutUnixNamoV2 string = "2006-01-02T15:04:05.000000"
)

// GetDayStartTimeStamp
//
//	@Description: 获取相应时区0点时间戳
//	@param now
//	@param offset
//	@return int64
func GetDayStartTimeStamp(timeStamp int64, offset int) int64 {
	loc := time.FixedZone("UTC", offset*60*60)
	curTime := time.Unix(timeStamp, 0).In(loc)
	return time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, loc).Unix()
}

// DayStartStr
//
//	@Description: 返回当时间字符串
//	@param timeStamp 时间戳
//	@param offset 时区
//	@param format 格式
//	@return string
func DayStartStr(timeStamp int64, offset int, format string) string {
	loc := time.FixedZone("UTC", offset*60*60)
	curTime := time.Unix(timeStamp, 0)
	tTime := curTime.In(loc)
	dayStart := time.Date(tTime.Year(), tTime.Month(), tTime.Day(), tTime.Hour(), tTime.Minute(), tTime.Second(), 0, loc)
	return dayStart.In(loc).Format(format)
}

// DayStrToTimeStamp 根据时间字符串转成时间戳
func DayStrToTimeStamp(timeStr string, offset int, format string) (timeStamp int64) {
	loc := time.FixedZone("LOC", offset*60*60)
	to, _ := time.ParseInLocation(format, timeStr, loc)
	timeStamp = to.Unix()
	return
}
