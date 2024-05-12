package converter

import (
	"github.com/arifullov/chat-server/internal/model"
	desc "github.com/arifullov/chat-server/pkg/chat_v1"
)

func ToChatFromDesc(chat *desc.CreateRequest) *model.Chat {
	return &model.Chat{
		Usernames: chat.GetUsernames(),
	}
}

func ToMessageFromDesc(chat *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		From:      chat.GetFrom(),
		Text:      chat.GetText(),
		ChatID:    chat.GetChatID(),
		Timestamp: chat.GetTimestamp().AsTime(),
	}
}
