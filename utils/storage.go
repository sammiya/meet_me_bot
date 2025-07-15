package utils

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// MeetingData は予定情報を格納する構造体
type MeetingData struct {
	Date                  *time.Time `json:"date"`
	PreNotificationSent   bool       `json:"pre_notification_sent"`
	LastNotificationDate  *time.Time `json:"last_notification_date"`
}

// LoadMeeting は JSONファイルから予定情報を読み込む
func LoadMeeting(filepath string) (*MeetingData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// ファイルが存在しない場合は空のデータを返す
			return &MeetingData{}, nil
		}
		return nil, err
	}

	var meeting MeetingData
	err = json.Unmarshal(data, &meeting)
	if err != nil {
		return nil, err
	}

	return &meeting, nil
}

// SaveMeeting は予定情報を JSONファイルに保存する
func SaveMeeting(filepath string, meeting *MeetingData) error {
	data, err := json.MarshalIndent(meeting, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// UpdateLastNotificationDate は最後の通知日を更新する
func UpdateLastNotificationDate(filepath string, date time.Time) error {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return err
	}

	meeting.LastNotificationDate = &date
	return SaveMeeting(filepath, meeting)
}

// GetLastNotificationDate は最後の通知日を取得する
func GetLastNotificationDate(filepath string) (*time.Time, error) {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return nil, err
	}

	return meeting.LastNotificationDate, nil
}

// SetNextWeekMeeting は1週間後の21:00に予定を設定する
func SetNextWeekMeeting(filepath string, baseDate time.Time) error {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return err
	}
	
	nextWeek := baseDate.AddDate(0, 0, 7)
	nextWeekAt21 := time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 21, 0, 0, 0, JST)
	
	meeting.Date = &nextWeekAt21
	meeting.PreNotificationSent = false
	// LastNotificationDateは既存の値を保持
	return SaveMeeting(filepath, meeting)
}

// AddMeeting は新しい予定を設定する
func AddMeeting(filepath string, date time.Time) error {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return err
	}
	
	meeting.Date = &date
	meeting.PreNotificationSent = false
	return SaveMeeting(filepath, meeting)
}

// ClearMeeting は予定情報をクリアする（最後の通知日は保持）
func ClearMeeting(filepath string) error {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return err
	}
	
	meeting.Date = nil
	meeting.PreNotificationSent = false
	return SaveMeeting(filepath, meeting)
}

// UpdatePreNotificationSent は1時間前通知済みフラグを更新する
func UpdatePreNotificationSent(filepath string) error {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return err
	}

	if meeting.Date == nil {
		return errors.New("no meeting scheduled")
	}

	meeting.PreNotificationSent = true
	return SaveMeeting(filepath, meeting)
}

// GetMeeting は予定情報を取得する（既存のAPIと互換性を保つ）
func GetMeeting(filepath string) (*time.Time, bool, error) {
	meeting, err := LoadMeeting(filepath)
	if err != nil {
		return nil, false, err
	}

	return meeting.Date, meeting.PreNotificationSent, nil
}