package utils

import (
	"os"
	"testing"
	"time"
)

func TestDatabase(t *testing.T) {
	// テスト用のデータベースファイルを作成
	testDB := "test.db"
	db, err := PrepareDB(testDB)
	if err != nil {
		t.Fatalf("Failed to prepare database: %v", err)
	}
	defer db.Close()
	defer os.Remove(testDB)

	// AddMeeting のテスト
	testDate := time.Date(2024, 3, 7, 21, 0, 0, 0, JST)
	err = AddMeeting(db, testDate)
	if err != nil {
		t.Errorf("Failed to add meeting: %v", err)
	}

	// GetMeeting のテスト
	meeting, preNotificationSent, err := GetMeeting(db)
	if err != nil {
		t.Errorf("Failed to get meeting: %v", err)
	}
	if meeting == nil {
		t.Error("Meeting is nil")
	} else {
		if !meeting.Equal(testDate) {
			t.Errorf("Expected date %v, got %v", testDate, meeting)
		}
	}
	if preNotificationSent {
		t.Error("Expected preNotificationSent to be false")
	}

	// UpdatePreNotificationSent のテスト
	err = UpdatePreNotificationSent(db)
	if err != nil {
		t.Errorf("Failed to update pre-notification sent: %v", err)
	}
	_, preNotificationSent, err = GetMeeting(db)
	if err != nil {
		t.Errorf("Failed to get meeting: %v", err)
	}
	if !preNotificationSent {
		t.Error("Expected preNotificationSent to be true")
	}

	// ClearMeetings のテスト
	err = ClearMeetings(db)
	if err != nil {
		t.Errorf("Failed to clear meetings: %v", err)
	}
	meeting, _, err = GetMeeting(db)
	if err != nil {
		t.Errorf("Failed to get meeting: %v", err)
	}
	if meeting != nil {
		t.Error("Expected meeting to be nil after clearing")
	}
}
