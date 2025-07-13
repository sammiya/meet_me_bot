package utils

import (
	"os"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	// テスト用の一時ファイル
	tempFile := "test_meeting.json"
	defer os.Remove(tempFile) // テスト後にファイルを削除

	// 初期状態のテスト
	meeting, preNotificationSent, err := GetMeeting(tempFile)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}
	if meeting != nil {
		t.Error("Expected no meeting initially")
	}
	if preNotificationSent {
		t.Error("Expected preNotificationSent to be false initially")
	}

	// 予定追加のテスト
	testDate := time.Date(2024, 3, 21, 21, 0, 0, 0, JST)
	err = AddMeeting(tempFile, testDate)
	if err != nil {
		t.Fatalf("AddMeeting failed: %v", err)
	}

	// 追加した予定の確認
	meeting, preNotificationSent, err = GetMeeting(tempFile)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}
	if meeting == nil {
		t.Fatal("Expected meeting to be set")
	}
	if !meeting.Equal(testDate) {
		t.Errorf("Expected meeting date %v, got %v", testDate, *meeting)
	}
	if preNotificationSent {
		t.Error("Expected preNotificationSent to be false")
	}

	// 通知済みフラグ更新のテスト
	err = UpdatePreNotificationSent(tempFile)
	if err != nil {
		t.Fatalf("UpdatePreNotificationSent failed: %v", err)
	}

	_, preNotificationSent, err = GetMeeting(tempFile)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}
	if !preNotificationSent {
		t.Error("Expected preNotificationSent to be true")
	}

	// 予定クリアのテスト
	err = ClearMeeting(tempFile)
	if err != nil {
		t.Fatalf("ClearMeeting failed: %v", err)
	}

	meeting, preNotificationSent, err = GetMeeting(tempFile)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}
	if meeting != nil {
		t.Error("Expected no meeting after clear")
	}
	if preNotificationSent {
		t.Error("Expected preNotificationSent to be false after clear")
	}
}