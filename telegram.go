package telegram

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

var Client *client.Client

type Config struct {
	ApiId      int32
	ApiHash    string
	DataFolder string
}

func Init(self *Config, opts ...client.Option) {
	authorizer := client.ClientAuthorizer(&client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join(self.DataFolder, ".tdlib", "database"),
		FilesDirectory:      filepath.Join(self.DataFolder, ".tdlib", "files"),
		UseFileDatabase:     false,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               self.ApiId,
		ApiHash:             self.ApiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		ApplicationVersion:  "1.0.0",
	})
	go client.CliInteractor(authorizer)
	var err error
	Client, err = client.NewClient(authorizer, opts...)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}
}

func GetHistory(chatId int64, FromMsgId int64, limit int32) (*client.Messages, error) {
	return Client.GetChatHistory(&client.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: FromMsgId,
		Limit:         limit,
	})
}

func GetMsg(chatId int64, id int64) (*client.Message, error) {
	return Client.GetMessage(&client.GetMessageRequest{
		ChatId:    chatId,
		MessageId: id,
	})
}

func SplitVideoMessage(msg *client.Message) (text string, fileId int32, ok bool) {
	var t *client.MessageVideo
	t, ok = msg.Content.(*client.MessageVideo)
	if !ok {
		return
	}

	text = t.Caption.Text
	fileId = t.Video.Video.Id

	return
}

func DownloadFile(fileId int32, sync bool) (*client.File, error) {
	return Client.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileId,
		Priority:    32,
		Synchronous: sync,
	})
}

func DownloadStat(fileId int32) (*client.File, error) {
	return Client.GetFile(&client.GetFileRequest{
		FileId: fileId,
	})
}

func GetListener() *client.Listener {
	return Client.GetListener()
}
