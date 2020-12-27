import * as crypto from 'crypto'
import {Chat, Message} from "../../model";
import {BaseRedis} from './base'
import { v4 as uuidv4 } from 'uuid';

export class ChatDbClient extends BaseRedis {
    constructor() {
        super()
    }

    private usersInfoSuffix = 'users'
    private userChatIdsSuffix = 'chats'
    private userSocketIdSuffix = 'socket'
    private messagesSuffix = 'messages'

    private addChatIdToUserChats(userId: string, chatId: string) {
        if (!this.client) {
            console.error("Client is undef for user: ", userId)
            return
        }

        this.client.sadd(userId + ':' + this.userChatIdsSuffix, chatId)
    }

    private addChatRecord(userIds: string[]): string {
        if (!this.client) {
            console.error("Client is undef (addChatRecord) ")
            return ''
        }

        const chatId = crypto
            .createHash('sha256')
            .update(userIds.concat(Date.now().toString()).join(' '))
            .digest('hex')

        this.client.sadd(chatId + ':' + this.usersInfoSuffix, userIds)
        return chatId
    }

    private addMessageIdToPool(messageId: string, chatId: string) {
        if (!this.client) {
            console.error("Client is undef (addMessageIdToPool) ")
            return ''
        }

        this.client.rpush(chatId + ':' + this.messagesSuffix, messageId)
    }

    private addMessageRecord(message: Message): string {
        if (!this.client) {
            console.error("Client is undef (addMessageRecord) ")
            return ''
        }

        const messageId = uuidv4()

        this.client.hmset(messageId,
            ['text', message.text,
                'status', message.status,
                'sender', message.sender,
                'recipient', message.recipient,
                'id', message.id,
                'date', message.date,
                'chatId', message.chatId])

        return messageId
    }

    private getUserChatIds(userId: string): Promise<string[]> {

        return new Promise<string[]>((resolve, reject) => {
            if (!this.client) {
                console.error("Client is undef (getUserChatIds) ")
                resolve([])
                return []
            }

            this.client.smembers(userId + ':' + this.userChatIdsSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting user chats: ", err)
                    resolve([])
                    return
                }
                resolve(res)
            })
        })
    }

    private getChatUsers(chatId: string): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            if (!this.client) {
                console.error("Client is undef (getChatUsers) ")
                resolve([])
                return
            }

            this.client.smembers(chatId + ':' + this.usersInfoSuffix,
                (err, res) => {
                if (err) {
                    console.log("Error getting chat users: ", err)
                    resolve([])
                    return
                }
                resolve(res)
            })
        })
    }

    private getChatMessages(chatId: string): Promise<Message[]> {
        return new Promise<Message[]>((resolve, reject) => {
            if (!this.client) {
                console.error("Client is undef (getChatMessages) ")
                resolve([])
                return
            }

            this.client.hmget(chatId + ':' + this.messagesSuffix,
                (err, res) => {
                    if (err) {
                        console.log("Error getting chat users: ", err)
                        resolve([])
                        return
                    }
                    console.log("Got messages: ", res);
                    
                    resolve(res.map(item => ({
                        text: item
                    } as Message)))
                })
        })
    }

    public registerUserAsOnline(userId: string, socketId: string) {
        if (!this.client) {
            console.error("Client is undef (registerUserAsOnline) ")
            return
        }
        this.client.set(userId + ':' + this.userSocketIdSuffix, socketId)
    }

    public unregisterUser(userId: string) {
        if (!this.client) {
            console.error("Client is undef (unregisterUser) ")
            return
        }
        this.client.set(userId + ':' + this.userSocketIdSuffix, '')
    }

    public getUserSocketId(userId: string): Promise<string | null> {
        return new Promise<string | null>((resolve, reject) => {
            if (!this.client) {
                console.error("Client is undef (getUserSocketId) ")
                resolve(null)
                return
            }

            this.client.get(userId + ':' + this.userSocketIdSuffix, (err, res) => {
                if (err) {
                    console.log("Error getting user socket io: ", err)
                    resolve(null)
                } else {
                    resolve(res)
                }
            })
        })
    }

    public createChat(userIds: string[]): string {
        const chatId = this.addChatRecord(userIds)

        userIds.forEach(id => this.addChatIdToUserChats(id, chatId))

        return chatId
    }

    public addMessageToChat(message: Message) {
        const id = this.addMessageRecord(message)
        this.addMessageIdToPool(id, message.chatId)
    }

    public async getUserChats(userId: string): Promise<Chat[]> {
        const chatIds = await this.getUserChatIds(userId)
        const chats = [] as Chat[]

        for (let i = 0; i < chatIds.length; i ++) {
            chats.push({
                id: chatIds[i],
                userIds: await this.getChatUsers(chatIds[i]),
                messages: await this.getChatMessages(chatIds[i])
            })
        }

        return chats
    }
}
