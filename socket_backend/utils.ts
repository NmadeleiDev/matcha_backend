import {Chat, Message, WSmessage} from "./model";
import CONSTANTS from "./config";
import * as storage from './storageStub'
import {Socket} from "socket.io";


export const sendToUser = (WStype: string, type: string, payload: object, userIds: string[] | null) => {
    if (WStype === CONSTANTS.WS.CHAT) {
        payload = payload as Chat
    } else if (WStype === CONSTANTS.WS.MESSAGE) {
        payload = payload as Message
    }

    if (userIds === null) {
        userIds = Array.from(storage.getUsers().values()).map((user) => user.id)
    }

    userIds.forEach(userId => {
        const user = storage.getUsers().get(userId);
        const wsMessage = new WSmessage({type, payload: payload});
        console.log("[sendToUser]", {userId, type, wsMessage});
        user && user.socket.emit(WStype, wsMessage.toString());
    })
};


export const addMessageToChat = (chatId: string, message: Message) => {
    const chat = storage.getChats().get(chatId);

    chat && chat.messages.push(message) && storage.getChats().set(chatId, chat);
};

export const findMessageInChat = (chatId: string, messageId: string): Message | null => {
    const chat = storage.getChats().get(chatId);
    if (!chat) {
        console.log(`Failed to find chat with id = ${chatId} to find message = ${messageId}`)
        return null
    }

    const message = chat.messages.find(item => item.id === messageId)
    if (!message) {
        console.log(`Failed to find message with id = ${messageId} in chat = ${chatId}`)
        return null
    }

    return message
}

export const updateMessageInChat = (chatId: string, payload: Message) => {
    const message = findMessageInChat(chatId, payload.id)

    if (!message) return

    message.text = payload.text
    message.status = payload.status
    // chat.messages.push(payload) && storage.getChats().set(chatId, chat);
};

export const deleteMessageFromChat = (chatId: string, payload: Message) => {
    const chat = storage.getChats().get(chatId);
    if (!chat) {
        console.log(`Failed to find chat with id = ${chatId} to find message = ${payload.id}`)
        return null
    }

    const idx = chat.messages.findIndex(item => item.id === payload.id)
    if (idx < 0) {
        console.log(`Failed to find message with id = ${payload.id}`)
        return
    }

    chat.messages.splice(idx, 1)
};


export const getUserId = (socket: Socket): string => {
    const user = storage.getUsers().get(socket.id);
    const id = user && user.id;
    if (id) return id;
    else throw new Error("User not found by socket id");
};

export const deleteUserFromChat = (userId: string) => {
    storage.getUsers().delete(userId);
};

export const getChatLength = (chatId: string): number => {
    const chat = storage.getChats().get(chatId);
    return chat ? chat.messages.length : 0;
};

export const mapChatsToIds = () => {
    return storage.getChats() && storage.getChats().size ? Array.from(storage.getChats().keys()) : [];
};

export const mapUsersToIds = () => {
    return storage.getUsers() && storage.getUsers().size ? Array.from(storage.getUsers().keys()) : [];
};

export const getAllUserChats = (userId: string | undefined): string[] => {
    console.log([...storage.getChats().values()]);
    [...storage.getChats().values()].forEach((el) => console.log(el));
    return userId && storage.getChats() && storage.getChats().size
        ? // if userId is in chat => return chat id
        [...storage.getChats().values()]
            .filter((chat) => !!chat.userIds.find((id) => id === userId))
            .map((chat) => chat.id)
        : [];
};
