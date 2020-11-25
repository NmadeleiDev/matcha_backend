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


func (c *chatsManager) SendMessageToChat(message model.Message) {
	log.Infof("Send message: %v, chatId: %v", message, message.ChatId)
	chat := c.FindChat(message.ChatId)
	if chat == nil {
		log.Errorf("Not sent message. Failed to find chat: %v; chats: %v", message.ChatId, c.Chats)
		return
	}

	if len(chat.UserIds) != 2 {
		log.Errorf("Not sent message. Chat is incorrect (not 2 members): %v", chat)
		return
	}
	message.State = 2
	chat.Messages = append(chat.Messages, message)

	wsMessage := model.SocketMessage{MessageType: 1, Payload: message}
	for _, id := range chat.UserIds {
		if recipient, exists := Clients[id]; recipient != nil && exists {
			recipient.ReadMessageChan <- wsMessage
		} else {
			log.Warnf("Message recipient is not online. Id = %v", message.Recipient)
		}
	}
}

func (c *chatsManager) UpdateMessageToChat(message model.Message) {
	chat := c.FindChat(message.ChatId)
	if chat == nil {
		log.Errorf("Not sent message. Failed to find chat: %v", message.ChatId)
		return
	}

	if len(chat.UserIds) != 2 {
		log.Errorf("Not sent message. Chat is incorrect (not 2 members): %v", chat)
		return
	}

	for i, _ := range chat.Messages {
		if chat.Messages[i].Id == message.Id {
			chat.Messages[i].Text = message.Text
		}
	}

	wsMessage := model.SocketMessage{MessageType: 2, Payload: message}
	for _, id := range chat.UserIds {
		if recipient, exists := Clients[id]; exists {
			recipient.ReadMessageChan <- wsMessage
		} else {
			log.Warnf("Message recipient is not online. Id = %v", message.Recipient)
		}
	}
}

func (c *chatsManager) DeleteMessageFromChat(message model.Message) {
	log.Infof("delete: %v", message)
	//chat := c.FindChat(message.ChatId)
	//if chat == nil {
	//	log.Errorf("Not sent message. Failed to find chat: %v", message.ChatId)
	//	return
	//}
	//
	//if len(chat.UserIds) != 2 {
	//	log.Errorf("Not sent message. Chat is incorrect (not 2 members): %v", chat)
	//	return
	//}
	//
	//for i, _ := range chat.Messages {
	//	if chat.Messages[i].Id == message.Id {
	//
	//	}
	//}
	//
	//wsMessage := model.SocketMessage{MessageType: 2, Payload: message}
	//for _, id := range chat.UserIds {
	//	if recipient, exists := Clients[id]; exists {
	//		recipient.ReadMessageChan <- wsMessage
	//	} else {
	//		log.Warnf("Message recipient is not online. Id = %v", message.Recipient)
	//	}
	//}
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
			log.Errorf("Recovered in CreateChat. Chat: %v; clients: %v; r = %v", chat, Clients, r)
		}
	}()
	chat.Id = hash.CalculateSha256(time.Now().String() + chat.UserIds[0])

	c.Chats = append(c.Chats, &chat)
	for _, id := range chat.UserIds {
		if client, exists := Clients[id]; client != nil && exists {
			client.ReadMessageChan <- model.SocketMessage{MessageType: 100, Payload: chat}
		} else {
			log.Infof("Client not exists or nil: %v", id)
		}
	}
	log.Infof("Created chat: %v", chat)
	return chat.Id
}

func (c *chatsManager) ConnectToChat(chatId string) {
	panic("implement me")
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
