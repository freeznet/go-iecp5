package asdu

import (
	"encoding/binary"
	"time"
)

// CP56Time2a , CP24Time2a, CP16Time2a
// |         Milliseconds(D7--D0)        | Milliseconds = 0-59999
// |         Milliseconds(D15--D8)       |
// | IV(D7)   RES1(D6)  Minutes(D5--D0)  | Minutes = 1-59, IV = invalid,0 = valid, 1 = invalid
// | SU(D7)   RES2(D6-D5)  Hours(D4--D0) | Hours = 0-23, SU = summer Time,0 = standard time, 1 = summer time,
// | DayOfWeek(D7--D5) DayOfMonth(D4--D0)| DayOfMonth = 1-31  DayOfWeek = 1-7
// | RES3(D7--D4)        Months(D3--D0)  | Months = 1-12
// | RES4(D7)            Year(D6--D0)    | Year = 0-99

// CP56Time2a time to CP56Time2a
func CP56Time2a(t time.Time, loc *time.Location) []byte {
	if loc == nil {
		loc = time.UTC
	}
	ts := t.In(loc)
	msec := ts.Nanosecond()/int(time.Millisecond) + ts.Second()*1000
	return []byte{byte(msec), byte(msec >> 8), byte(ts.Minute()), byte(ts.Hour()),
		byte(ts.Weekday()<<5) | byte(ts.Day()), byte(ts.Month()), byte(ts.Year() - 2000)}
}

// 7个八位位组二进制时间，建议所有时标采用UTC
// The year is assumed to be in the 20th century.
// See IEC 60870-5-4 § 6.8 and IEC 60870-5-101 second edition § 7.2.6.18.
func ParseCP56Time2a(bytes []byte, loc *time.Location) time.Time {
	if len(bytes) < 7 || bytes[2]&0x80 == 0x80 {
		return time.Time{}
	}

	x := int(binary.LittleEndian.Uint16(bytes))
	msec := x % 1000
	sec := x / 1000
	min := int(bytes[2] & 0x3f)
	hour := int(bytes[3] & 0x1f)
	day := int(bytes[4] & 0x1f)
	month := time.Month(bytes[5] & 0x0f)
	year := 2000 + int(bytes[6]&0x7f)

	nsec := msec * int(time.Millisecond)
	if loc == nil {
		loc = time.UTC
	}
	return time.Date(year, month, day, hour, min, sec, nsec, loc)
}

func CP24Time2a(t time.Time, loc *time.Location) []byte {
	if loc == nil {
		loc = time.UTC
	}
	ts := t.In(loc)
	msec := ts.Nanosecond()/int(time.Millisecond) + ts.Second()*1000
	return []byte{byte(msec), byte(msec >> 8), byte(ts.Minute())}
}

// 3个八位位组二进制时间，建议所有时标采用UTC
// 读3字节,返回一个值，当无效时返回空值
// The moment is assumed to be in the recent present.
// See IEC 60870-5-4 § 6.8 and IEC 60870-5-101 second edition § 7.2.6.19.
func ParseCP24Time2a(bytes []byte, loc *time.Location) time.Time {
	if len(bytes) < 3 || bytes[2]&0x80 == 0x80 {
		return time.Time{}
	}
	x := int(binary.LittleEndian.Uint16(bytes))
	msec := x % 1000
	sec := (x / 1000)
	min := int(bytes[2] & 0x3f)
	now := time.Now()
	year, month, day := now.Date()
	hour, _, _ := now.Clock()

	nsec := msec * int(time.Millisecond)
	if loc == nil {
		loc = time.UTC
	}
	val := time.Date(year, month, day, hour, min, sec, nsec, loc)

	////5 minute rounding - 55 minute span
	//if min > currentMin+5 {
	//	val = val.Add(-time.Hour)
	//}

	return val
}

func CP16Time2a(msec uint16) []byte {
	return []byte{byte(msec), byte(msec >> 8)}
}

func ParseCP16Time2a(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}
