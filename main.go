package main

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"strings"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/joho/godotenv"
	"io"
	"github.com/google/uuid"
	"net/http"
	"os"
	"fmt"
	"bytes"
)

func sendHello(bot *telego.Bot, update telego.Update) {
	bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		fmt.Sprintf("Hello, %s, nice to meet you, send me a file and i will send you a link to my dufs with your file", update.Message.From.FirstName),
	))
}

func asyncFileHandler(bot *telego.Bot, update telego.Update) {
	if update.Message.Document == nil && update.Message.Audio == nil{
		bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			"There is no file in this message",
		))
		return
	}

	var fileid string
	var got_file *telego.File
	var bts []byte

	if update.Message.Document != nil {
		fileid = fmt.Sprintf("%s.%s", uuid.NewString(), strings.Split(update.Message.Document.FileName, ".")[len(strings.Split(update.Message.Document.FileName, ".")) - 1])
		got_file, _ = bot.GetFile(&telego.GetFileParams{FileID: update.Message.Document.FileID})
		bts, _ = tu.DownloadFile(bot.FileDownloadURL(got_file.FilePath))
	} else {
		fileid = fmt.Sprintf("%s.%s", uuid.NewString(), strings.Split(update.Message.Audio.FileName, ".")[len(strings.Split(update.Message.Audio.FileName, ".")) - 1])
		got_file, _ = bot.GetFile(&telego.GetFileParams{FileID: update.Message.Audio.FileID})
		bts, _ = tu.DownloadFile(bot.FileDownloadURL(got_file.FilePath))
	}

	body := &bytes.Buffer{}
	reader := bytes.NewReader(bts)
	writer := io.Writer(body)
	_, _ = io.Copy(writer, reader)
	
	dufsURL := os.Getenv("DUFS_URL")
	request, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", dufsURL, fileid), body)
	client := &http.Client{}
	client.Do(request)
	bot.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Text: fmt.Sprintf("Your file is here: %s/%s", dufsURL, fileid),
	})

}

func handleFile(bot* telego.Bot, update telego.Update) {
	go asyncFileHandler(bot, update)
}

func main() {
	_ = godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	bot, _ := telego.NewBot(token, telego.WithDefaultDebugLogger())

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(sendHello, th.CommandEqual("start"))
	bh.Handle(handleFile,th.AnyMessage())

	bh.Start()
}
