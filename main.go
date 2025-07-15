package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"meet_me_bot/utils"
)

// JSON データファイルの名前
var dataFile = "meet_me_bot.json"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	channelID := os.Getenv("DISCORD_CHANNEL")
	announcerID := os.Getenv("DISCORD_MEETING_ANNOUNCER")

	// 上記値のどれかが空文字列の場合はエラーを出力して終了
	if token == "" || channelID == "" || announcerID == "" {
		log.Fatal("DISCORD_TOKEN, DISCORD_CHANNEL, DISCORD_MEETING_ANNOUNCER must be set")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	dg.AddHandler(messageCreate(channelID, announcerID, dataFile))

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	// ticker := time.NewTicker(time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			scheduledMessage(dg, channelID, dataFile)
		}
	}()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(channelID string, announcerID string, dataFile string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != channelID {
			return
		}

		if utils.IsAskingNextDate(m.Content) {
			meeting, _, err := utils.GetMeeting(dataFile)

			if err != nil {
				log.Println("Error getting meeting: ", err)
				return
			}

			if meeting == nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "次回予定は未定だよ")
				if err != nil {
					log.Println("Error sending message: ", err)
				}
				return
			}

			// meeting は *time.Time 型
			_, err = s.ChannelMessageSend(m.ChannelID, "次回予定: "+utils.FormatDate(*meeting))
			if err != nil {
				log.Println("Error sending message: ", err)
			}
		}

		if announcerID != m.Author.ID {
			return
		}

		// 無言メンションのチェック（メンションしてくれるだけでテキストなし）
		if len(m.Mentions) > 0 && m.Content == "<@"+s.State.User.ID+">" {
			lastNotificationDate, err := utils.GetLastNotificationDate(dataFile)
			if err != nil {
				log.Println("Error getting last notification date: ", err)
				return
			}

			if lastNotificationDate == nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "最後の通知日が見つからないよ")
				if err != nil {
					log.Println("Error sending message: ", err)
				}
				return
			}

			err = utils.SetNextWeekMeeting(dataFile, *lastNotificationDate)
			if err != nil {
				log.Println("Error setting next week meeting: ", err)
				return
			}

			// 新しい予定を取得して通知
			nextDate, _, err := utils.GetMeeting(dataFile)
			if err != nil {
				log.Println("Error getting meeting: ", err)
				return
			}

			if nextDate != nil {
				_, err = s.ChannelMessageSend(m.ChannelID, "次回予定: "+utils.FormatDate(*nextDate))
				if err != nil {
					log.Println("Error sending message: ", err)
				}
			}
			return
		}

		parsedDate, err := utils.ParseNextDate(time.Now().In(utils.JST), m.Content)

		if err != nil {
			log.Println("Error parsing date: ", err)
			return
		}

		if parsedDate != nil {
			err := utils.AddMeeting(dataFile, *parsedDate)

			if err != nil {
				err = s.MessageReactionAdd(m.ChannelID, m.ID, "👎")

				if err != nil {
					log.Println("Error adding reaction 👎: ", err)
				}

				return
			}
			err = s.MessageReactionAdd(m.ChannelID, m.ID, "👍")

			if err != nil {
				log.Println("Error adding reaction: ", err)
			}
		}
	}
}

func scheduledMessage(dg *discordgo.Session, channelID string, dataFile string) {
	fmt.Println("Scheduled message")

	nextTime, preNotificationSent, err := utils.GetMeeting(dataFile)

	if err != nil {
		log.Println("Error getting meeting: ", err)
		return
	}

	if nextTime == nil {
		return
	}

	remainingSeconds := time.Until(*nextTime).Seconds()

	if !preNotificationSent && 3300 < remainingSeconds && remainingSeconds < 3660 {
		remainingForPrenotification := remainingSeconds - 3600
		time.Sleep(time.Duration(remainingForPrenotification) * time.Second)

		dg.ChannelMessageSend(channelID, "1時間前だよ")

		utils.UpdatePreNotificationSent(dataFile)
	}

	if remainingSeconds < 60 {
		time.Sleep(time.Duration(remainingSeconds) * time.Second)

		dg.ChannelMessageSend(channelID, "はじまるよ")

		// 最後の通知日を保存
		utils.UpdateLastNotificationDate(dataFile, *nextTime)

		utils.ClearMeeting(dataFile)
	}

}
