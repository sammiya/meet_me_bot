package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestParseNextDate(t *testing.T) {
    // テストケースを定義
    testCases := []struct {
        name     string
        text     string
        expected *time.Time
    }{
        {
            name: "Case 1",
            text: "次回は 03/07 ね。",
            expected: func() *time.Time {
                d := time.Date(2024, 3, 7, 21, 0, 0, 0, JST)
                return &d
            }(),
        },
        {
            name: "Case 2",
            text: "次回12/05",
            expected: func() *time.Time {
				d := time.Date(2024, 12, 5, 21, 0, 0, 0, JST)
				return &d
			}(),
        },
        {
            name: "Case 3",
            text: "次回 05/13（木）",
            expected: func() *time.Time {
				d := time.Date(2024, 5, 13, 21, 0, 0, 0, JST)
				return &d
			}(),
        },
        {
            name:     "Case 4",
            text:     "foo",
            expected: nil, // 日付が含まれない場合はnilを返す
        },
    }

    // テストを実行
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, _ := ParseNextDate(
				time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC),
				tc.text,
			)
            if !reflect.DeepEqual(result, tc.expected) {
                t.Errorf("Expected %v, but got %v", tc.expected, result)
            }
        })
    }
}

func TestFormatDate(t *testing.T) {
    // テストケースを定義
    testCases := []struct {
        name     string
        date     time.Time
        expected string
    }{
        {
            name:     "Case 1",
            date:     time.Date(2024, 3, 7, 21, 0, 0, 0, JST),
            expected: "3/7（木）",
        },
        {
            name:     "Case 2",
            date:     time.Date(2024, 12, 5, 21, 0, 0, 0, JST),
            expected: "12/5（木）",
        },
        {
            name:     "Case 3",
            date:     time.Date(2024, 5, 13, 21, 0, 0, 0, JST),
            expected: "5/13（月）",
        },
    }

    // テストを実行
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := FormatDate(tc.date)
            if result != tc.expected {
                t.Errorf("Expected %s, but got %s", tc.expected, result)
            }
        })
    }
}
