import {Socket} from "socket.io";
import ChatsDb from '../db/redis'
import {Chat, Message} from "../model/model";

const Clients = new Map() as Map<string, Socket>

export function registerSocketConn(socket: Socket, clientId: string) {
    ChatsDb.registerUserAsOnline(clientId, socket.id)

    Clients.set(clientId, socket)

    ChatsDb.getUserChats(clientId)
        .then(chats => {
            socket.send(chats)

            chats.forEach(chat => {
                socket.join(chat.id) // подключаем пользователя к его чатам
            })
        })

    socket.on('message', (msg) => {
        console.log('message: ', msg);

        const message = msg as Message

        message.state = 2

        ChatsDb.addMessageToChat(message)

        socket.to(message.chatId).emit('message', message) // под каждый чат создается комната, к которой подключены оналайн пользователи
    });
}

export function sendChatCreationMessage(chat: Chat) {
    chat.userIds.forEach(id => {
        const client = Clients.get(id)

        if (client) {
            client.emit('chat', chat)
        } else {
            console.log("Client is not online: ", id)
        }
    })
}