package notifyserv

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/domain/entities"
	"time"
)

type telegramBot struct {
	cfg *config.Config
	bot *tgbotapi.BotAPI
}

func NewTelegramBot(cfg *config.Config) NotifyService {
	bot, err := tgbotapi.NewBotAPI(cfg.Notify.Token)
	if err != nil {
		log.Panic(err)
	}

	return &telegramBot{
		cfg: cfg,
		bot: bot,
	}
}

func (t telegramBot) SendMessageToTelegram(trans *entities.TransactionLog, reason string) error {
	log.Info("Sending message to telegram")

	transactionHeader := "Transaction is failed"
	if trans.TransactionStatus == entities.TX_SUCCESS {
		transactionHeader = "Transaction is finished"
	}

	date := time.Now().Format("2006-01-02 15:04:05")
	detailURL := fmt.Sprintf("http://localhost:5020/api/v1/transaction/detail/%s", trans.ID.Hex())
	text := fmt.Sprintf("%s\n\nOrderID: %s\nStatus: %v\nReason: %s\nDate: %s\n\n %s",
		transactionHeader, trans.OrderID, trans.TransactionStatus, reason, date, detailURL)

	msg := tgbotapi.NewMessage(t.cfg.Notify.GroupAlertID, text)
	msg.ParseMode = "markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
