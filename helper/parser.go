package helper

import "time"

func ParseTime(stringTime string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, stringTime)

	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

func ParseTimeString(date time.Time) string {
	layout := "2006-01-02"
	return date.Format(layout)
}
