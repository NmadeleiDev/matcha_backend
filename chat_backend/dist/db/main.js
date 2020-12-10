"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const redis = __importStar(require("redis"));
const crypto = __importStar(require("crypto"));
const uuid_1 = require("uuid");
class ChatDbClient {
    constructor() {
        this.usersInfoSuffix = 'users';
        this.userChatIdsSuffix = 'chats';
        this.userSocketIdSuffix = 'socket';
        this.messagesSuffix = 'messages';
        this.initConnection();
    }
    initConnection() {
        const host = 'redis';
        const port = (process.env.REDIS_PORT && process.env.REDIS_PORT.length > 0)
            ? process.env.REDIS_PORT : '6379';
        const user = (process.env.REDIS_USER && process.env.REDIS_USER.length > 0)
            ? process.env.REDIS_USER : 'chat_user';
        const password = (process.env.REDIS_PASSWORD && process.env.REDIS_PASSWORD.length > 0)
            ? process.env.REDIS_PASSWORD : 'ffa9203c493aa99';
        this._client = redis.createClient({
            url: `redis://${user}:${password}@${host}:${port}`
        });
        this._client.on("error", function (error) {
            console.error(error);
        });
    }
    addChatIdToUserChats(userId, chatId) {
        this._client.sadd(userId + ':' + this.userChatIdsSuffix, chatId);
    }
    addChatRecord(userIds) {
        const chatId = crypto
            .createHash('sha256')
            .update(userIds.concat(Date.now().toString()).join(' '))
            .digest('hex');
        this._client.sadd(chatId + ':' + this.usersInfoSuffix, userIds);
        return chatId;
    }
    addMessageIdToPool(messageId, chatId) {
        this._client.rpush(chatId + ':' + this.messagesSuffix, messageId);
    }
    addMessageRecord(message) {
        const messageId = uuid_1.v4();
        this._client.hmset(messageId, ['text', message.text,
            'state', message.state,
            'sender', message.sender,
            'recipient', message.recipient,
            'id', message.id,
            'date', message.date,
            'chatId', message.chatId]);
        return messageId;
    }
    getUserChatIds(userId) {
        return new Promise((resolve, reject) => {
            this._client.smembers(userId + ':' + this.userChatIdsSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting user chats: ", err);
                    resolve([]);
                    return;
                }
                resolve(res);
            });
        });
    }
    getChatUsers(chatId) {
        return new Promise((resolve, reject) => {
            this._client.smembers(chatId + ':' + this.usersInfoSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting chat users: ", err);
                    resolve([]);
                    return;
                }
                resolve(res);
            });
        });
    }
    getChatMessages(chatId) {
        return new Promise((resolve, reject) => {
            this._client.hmget(chatId + ':' + this.messagesSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting chat users: ", err);
                    resolve([]);
                    return;
                }
                resolve(res);
            });
        });
    }
    registerUserAsOnline(userId, socketId) {
        this._client.set(userId + ':' + this.userSocketIdSuffix, socketId);
    }
    unregisterUser(userId) {
        this._client.set(userId + ':' + this.userSocketIdSuffix, undefined);
    }
    getUserSocketId(userId) {
        return new Promise((resolve, reject) => {
            this._client.get(userId + ':' + this.userSocketIdSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting user socket io: ", err);
                    resolve(null);
                }
                else {
                    resolve(res);
                }
            });
        });
    }
    createChat(userIds) {
        const chatId = this.addChatRecord(userIds);
        userIds.forEach(id => this.addChatIdToUserChats(id, chatId));
        return chatId;
    }
    addMessageToChat(message) {
        const id = this.addMessageRecord(message);
        this.addMessageIdToPool(id, message.chatId);
    }
    async getUserChats(userId) {
        const chatIds = await this.getUserChatIds(userId);
        const chats = [];
        for (let i = 0; i < chatIds.length; i++) {
            chats.push({
                id: chatIds[i],
                userIds: await this.getChatUsers(chatIds[i]),
                messages: await this.getChatMessages(chatIds[i])
            });
        }
        return chats;
    }
}
exports.default = new ChatDbClient();
//# sourceMappingURL=main.js.map