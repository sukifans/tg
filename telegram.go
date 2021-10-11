package telegram

import (
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

var c *client.Client

func Init(apiId int32, apiHash string) {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)
	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}
	var err error
	c, err = client.NewClient(authorizer, client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 0,
	}))
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}
}

func GetHistory(chatId int64, FromMsgId int64, limit int32) (*client.Messages, error) {
	return c.GetChatHistory(&client.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: FromMsgId,
		Limit:         limit,
	})
}

func GetMsg(chatId int64, id int64) (*client.Message, error) {
	return c.GetMessage(&client.GetMessageRequest{
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

func DownloadFile(fileId int32) (*client.File, error) {
	return c.DownloadFile(&client.DownloadFileRequest{
		FileId:   fileId,
		Priority: 32,
	})
}

func DownloadStat(fileId int32) (*client.File, error) {
	return c.GetFile(&client.GetFileRequest{
		FileId: fileId,
	})
}

func GetListener() *client.Listener {
	return c.GetListener()
}
