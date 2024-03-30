package utils

import (
	"fmt"
	"regexp"
	"time"
)

var JST = time.FixedZone("Asia/Tokyo", 9*60*60)

var WEEKDAYS = []string{
	"日", "月", "火", "水", "木", "金", "土",
}

func ParseNextDate(now time.Time, text string) (*time.Time, error) {
	re := regexp.MustCompile(`次回[^\d]*(\d{1,2})/(\d{1,2})`)
	match := re.FindStringSubmatch(text)
	if match == nil {
		return nil, nil
	}

	dateStr := fmt.Sprintf("%d/%02s/%02s 21:00", now.Year(), match[1], match[2])

	nextDate, err := time.ParseInLocation("2006/01/02 15:04", dateStr, JST)
	if err != nil {
		return nil, err
	}

	return &nextDate, nil
}

func IsAskingNextDate(text string) bool {
	pattern := `次.*いつ`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

// FormatDate は日付をフォーマットして返す。11/8（水） のように
func FormatDate(date time.Time) string {
	return fmt.Sprintf("%d/%d（%s）", date.Month(), date.Day(), WEEKDAYS[date.Weekday()])
}
