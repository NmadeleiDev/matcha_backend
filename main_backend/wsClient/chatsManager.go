package wsClient

import (
	"time"

	"backend/hash"
	"backend/model"
	"backend/utils"

	log "github.com/sirupsen/logrus"
)

type chatsManager struct {
	Chats []*model.Chat
}

func (c *chatsManager) FindChat(chatId string) *model.Chat {
	for _, chat := range c.Chats {
		if chat.Id == chatId {
			return chat
		}
	}
	return nil
}

func (c *chatsManager) GetUserChats(userId string) (result []*model.Chat) {
	for _, chat := range c.Chats {
		if utils.DoesArrayContain(chat.UserIds, userId) {
			result = append(result, chat)
		}
	}
	return result
}

func (c *chatsManager) CreateChat(chat model.Chat) string {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered in CreateChat. Chat: %v; clients: %v", chat, Clients)
		}
	}()
	chat.Id = hash.CalculateSha256(time.Now().String() + chat.UserIds[0])

	c.Chats = append(c.Chats, &chat)
	for _, id := range chat.UserIds {
		if client, exists := Clients[id]; exists {
			client.ReadMessageChan <- model.SocketMessage{MessageType: 100, ToChat: chat.Id, Payload: chat}
		} else {
			log.Infof("Client not exists: %v", id)
		}
	}
	log.Infof("Created chat: %v", chat)
	return chat.Id
}

func (c *chatsManager) ConnectToChat(chatId string) {
	panic("implement me")
}

func (c *chatsManager) SendMessageToChat(chatId string, message model.Message) {
	chat := c.FindChat(chatId)
	if chat == nil {
		log.Errorf("Not sent message. Failed to find chat: %v", chatId)
		return
	}

	if len(chat.UserIds) != 2 {
		log.Errorf("Not sent message. Chat is incorrect (not 2 members): %v", chat)
		return
	}

	chat.Messages = append(chat.Messages, message)

	wsMessage := model.SocketMessage{MessageType: 1, Payload: message, ToChat: chat.Id}
	for _, id := range chat.UserIds {
		if recipient, exists := Clients[id]; exists {
			recipient.ReadMessageChan <- wsMessage
		} else {
			log.Warnf("Message recipient is not online. Id = %v", message.Recipient)
		}
	}

}

func (c *chatsManager) AddUserToChat(userId string, destChat model.Chat) {
	for _, chat := range c.Chats {
		if chat.Id == destChat.Id {
			chat.UserIds = append(chat.UserIds, userId)
		}
	}
}

var manager chatsManager

func GetManager() model.WsDataManager {
	return &manager
}
