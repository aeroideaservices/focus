package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type StringInt int

func (st *StringInt) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}

	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err

		}
		*st = StringInt(i)
	}

	return nil
}

type Filesize int64

const (
	Bi Filesize = 1 << (iota * 10)
	KiB
	MiB
	GiB
	TiB
)

var sizesNames = []string{"KB", "MB", "GB", "TB"}

func (size Filesize) String() string {
	const unit = 1000 // 1000
	if size < unit {  // если размер меньше 1000, возвращаем в байтах
		return fmt.Sprintf("%d B", size)
	}

	exp := 0
	n := size >> 10
	for ; n >= unit && exp < len(sizesNames)-1; n = n >> 10 {
		exp++
	}
	n = 1 << (10 * (exp + 1))
	return fmt.Sprintf("%.1f %s", float64(size)/float64(n), sizesNames[exp])
}

func (size Filesize) MarshalJSON() ([]byte, error) {
	return json.Marshal(size.String())
}

type Time time.Time

const (
	TimeFormat = "2006-01-02 15:04:05"
)

func Now() Time {
	return Time(time.Now())
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}
