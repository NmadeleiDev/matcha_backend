import {Socket} from "socket.io";
import CONSTANTS from './config'
import {Chat, Message, WSmessage} from "./model";
import * as utils from './utils'
import * as storage from './storageStub'
import {sendToUser} from "./utils";
import {Express} from "express";
import * as queryString from "querystring";
import {MongoManager} from './db/mongo/mongo'
import {notificationsClient} from "./db/redis/notifications";

export function addSocketHandlers(userId: string, socket: Socket) {
    // debug only
    socket.onAny((event, ...args) => console.log("On any log: ", event, ...args));

    socket.on(CONSTANTS.WS.MESSAGE, (body) => {
        const json = JSON.parse(body) as WSmessage;
        const payload = json.payload as Message;
        switch (json.type) {
            case CONSTANTS.MESSAGE_TYPES.NEW_MESSAGE:
                // if first message in chat - send NEW_CHAT message first
                // save chat in db
                // add message to db
                // send to other NEW_MESSAGE
                // ---- or maybe just send DELIVERED_MESSAGE

                console.log(`[${json.type}]`)
                utils.addMessageToChat(payload.chatId, payload)

                payload.status = CONSTANTS.MESSAGE_STATUS.STATUS_DELIVERED

                utils.sendToUser(
                    CONSTANTS.WS.MESSAGE,
                    CONSTANTS.MESSAGE_TYPES.NEW_MESSAGE,
                    payload,
                    [payload.recipient]
                )

                utils.sendToUser(
                    CONSTANTS.WS.MESSAGE,
                    CONSTANTS.MESSAGE_TYPES.UPDATE_MESSAGE,
                    payload,
                    [payload.sender]
                )

                return;
            case CONSTANTS.MESSAGE_TYPES.UPDATE_MESSAGE:
                console.log(`[${json.type}]`);

                utils.updateMessageInChat(payload.chatId, payload);

                utils.sendToUser(
                    CONSTANTS.WS.MESSAGE,
                    CONSTANTS.MESSAGE_TYPES.UPDATE_MESSAGE,
                    payload,
                    [payload.recipient, payload.sender]
                );

                return;
            case CONSTANTS.MESSAGE_TYPES.DELETE_MESSAGE:
                console.log(`[${json.type}]`);

                utils.deleteMessageFromChat(payload.chatId, payload);

                utils.sendToUser(
                    CONSTANTS.WS.MESSAGE,
                    CONSTANTS.MESSAGE_TYPES.UPDATE_MESSAGE,
                    payload,
                    [payload.recipient]
                );
                return;
            default:
                console.log("Unknown message type: ", json.type)
                return;
        }
    });

    socket.on(CONSTANTS.WS.CHAT, (body) => {
        const json = JSON.parse(body);
        const payload = json.payload as Chat;
        console.log("Chat type payload: ", json);
        switch (json.type) {
            case CONSTANTS.CHAT_TYPES.NEW_CHAT:
                const newChat = new Chat({userIds: payload.userIds});
                // save chat in memory, wait for the first message
                console.log("[NEW_CHAT]");
                const chatExists = [...storage.getChats().values()].find(
                    (chat) =>
                        chat.userIds.includes(newChat.userIds[0]) &&
                        chat.userIds.includes(newChat.userIds[1])
                );

                console.log(chatExists);
                if (chatExists) {
                    console.log("Chat to create already exists: ", chatExists)
                    sendToUser(CONSTANTS.WS.CHAT, CONSTANTS.CHAT_TYPES.NEW_CHAT, chatExists, chatExists.userIds)
                } else {
                    storage.getChats().set(newChat.id, newChat);
                    sendToUser(CONSTANTS.WS.CHAT, CONSTANTS.CHAT_TYPES.NEW_CHAT, newChat, newChat.userIds)
                }
                return
            case CONSTANTS.CHAT_TYPES.DELETE_CHAT:
                // send to both DELETE_CHAT message
                // delete chat from db
                console.log("[DELETE_CHAT]");

                const chatToDelete = storage.getChats().get(payload.id);
                if (!chatToDelete) {
                    console.log("Error find chat to delete: ", chatToDelete)
                    return;
                }

                utils.sendToUser(
                    CONSTANTS.WS.CHAT,
                    CONSTANTS.CHAT_TYPES.DELETE_CHAT,
                    chatToDelete,
                    chatToDelete.userIds
                );

                storage.getChats().delete(payload.id)
                return;
            default:
                console.log("Unknown chat type: ", json.type)
                return;
        }
    });

    socket.on("error", (e) => {
        console.log(e);
    });

    socket.on("disconnect", (reason) => {
        console.log("disconnect: " + reason);
        utils.setOnlineState(userId, false)
        notificationsClient.unsubscribe(userId)
        try {
            // utils.deleteUserFromChat(socket.id);
            console.log("user disconnected: " + socket.id);
        } catch (e) {
            console.log(e);
        }
    });
}

export function addServerHandlers(app: Express) {
    app.get("/test", ((req, res) => {
        const parsed = queryString.parse(req.url.split("?")[1]);
        res.json({status: true, data: parsed})
    }))

    app.get("/chats", ((req, res) => {
        // const parsed = queryString.parse(req.url.split("?")[1]);
        const id = req.query.id;
        console.log("Parsed get string: ", id)

        const userChats = Array.from(storage.getChats().values())
            .filter(item => item.userIds.includes(id as string))

        console.log("Found user chats: ", userChats)

        res.json({status: true, data: userChats})
    }))

    app.post("/chat", ((req, res) => {
        const newChat = new Chat({userIds: (req.body as Chat).userIds as string[]})

        const chatExists = [...storage.getChats().values()].find(
            (chat) =>
                chat &&
                chat.userIds.includes(newChat.userIds[0]) &&
                chat.userIds.includes(newChat.userIds[1])
        );

        console.log(chatExists);
        if (chatExists) {
            console.log("Chat to create already exists: ", chatExists)
            sendToUser(CONSTANTS.WS.CHAT, CONSTANTS.CHAT_TYPES.NEW_CHAT, chatExists, chatExists.userIds)
            res.json({status: true, data: chatExists})
        } else {
            storage.getChats().set(newChat.id, newChat);
            sendToUser(CONSTANTS.WS.CHAT, CONSTANTS.CHAT_TYPES.NEW_CHAT, newChat, newChat.userIds)
            res.json({status: true, data: newChat})
        }
    }))
}