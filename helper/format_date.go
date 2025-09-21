package helper

import "time"

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDate(dateString string) (time.Time, error) {
	layouts := []string{"2006-01-02", "02-01-2006"}
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}
	return parsedTime, err
}
