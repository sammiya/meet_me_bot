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

// ClearMeeting は予定情報をクリアする
func ClearMeeting(filepath string) error {
	emptyMeeting := &MeetingData{}
	return SaveMeeting(filepath, emptyMeeting)
}

// AddMeeting は新しい予定を設定する
func AddMeeting(filepath string, date time.Time) error {
	meeting := &MeetingData{
		Date:                &date,
		PreNotificationSent: false,
	}
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