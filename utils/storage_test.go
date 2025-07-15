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

	// 最後の通知日設定のテスト
	lastNotificationDate := time.Date(2024, 3, 14, 21, 0, 0, 0, JST)
	err = UpdateLastNotificationDate(tempFile, lastNotificationDate)
	if err != nil {
		t.Fatalf("UpdateLastNotificationDate failed: %v", err)
	}

	// 最後の通知日取得のテスト
	retrievedDate, err := GetLastNotificationDate(tempFile)
	if err != nil {
		t.Fatalf("GetLastNotificationDate failed: %v", err)
	}
	if retrievedDate == nil {
		t.Fatal("Expected last notification date to be set")
	}
	if !retrievedDate.Equal(lastNotificationDate) {
		t.Errorf("Expected last notification date %v, got %v", lastNotificationDate, *retrievedDate)
	}

	// 1週間後の予定設定のテスト
	err = SetNextWeekMeeting(tempFile, lastNotificationDate)
	if err != nil {
		t.Fatalf("SetNextWeekMeeting failed: %v", err)
	}

	meeting, _, err = GetMeeting(tempFile)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}
	if meeting == nil {
		t.Fatal("Expected meeting to be set")
	}

	// 1週間後の21:00になっているか確認
	expectedDate := time.Date(2024, 3, 21, 21, 0, 0, 0, JST)
	if !meeting.Equal(expectedDate) {
		t.Errorf("Expected next week meeting %v, got %v", expectedDate, *meeting)
	}

	// 予定クリアのテスト（最後の通知日は保持）
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

	// 最後の通知日が保持されているか確認
	retrievedDate, err = GetLastNotificationDate(tempFile)
	if err != nil {
		t.Fatalf("GetLastNotificationDate failed: %v", err)
	}
	if retrievedDate == nil {
		t.Error("Expected last notification date to be preserved after clear")
	} else if !retrievedDate.Equal(lastNotificationDate) {
		t.Errorf("Expected preserved last notification date %v, got %v", lastNotificationDate, *retrievedDate)
	}
}