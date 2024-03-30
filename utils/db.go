package utils

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// PrepareDB は、指定したファイルパスに SQLite データベースを作成し、データベースへの接続を返す
func PrepareDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS meetings (
			date TEXT NOT NULL,
			pre_notification_sent INTEGER NOT NULL
		)
		`)

	if err != nil {
		return nil, err
	}

	return db, nil
}

// ClearMeetings は meetings テーブルのレコードを全て削除する
func ClearMeetings(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM meetings")
	return err
}

// AddMeeting は meetings テーブルに新しいミーティングを追加する（既存のミーティングは削除される。また、pre_notification_sent は 0 で初期化される）
func AddMeeting(db *sql.DB, date time.Time) error {
	err := ClearMeetings(db)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO meetings (date, pre_notification_sent) VALUES (?, ?)", date.Format("2006-01-02 15:04:05"), 0)
	return err
}

// UpdatePreNotificationSent は meetings テーブルの pre_notification_sent を 1 に更新する
func UpdatePreNotificationSent(db *sql.DB) error {
	_, err := db.Exec("UPDATE meetings SET pre_notification_sent = 1")
	return err
}

// GetMeeting は meetings テーブルの最初のレコードを取得する（レコードが存在しない場合は nil を返す）
func GetMeeting(db *sql.DB) (*time.Time, bool, error) {
	var dateStr string
	var preNotificationSent int
	err := db.QueryRow("SELECT date, pre_notification_sent FROM meetings").Scan(&dateStr, &preNotificationSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	date, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, JST)
	if err != nil {
		return nil, false, err
	}

	return &date, preNotificationSent == 1, nil
}
