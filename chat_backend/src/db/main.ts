import * as redis from 'redis'
import * as crypto from 'crypto'
import {Chat, Message} from "../model/model";
import { v4 as uuidv4 } from 'uuid';

class ChatDbClient {
    constructor() {
        this.initConnection()
    }

    private _client: redis.RedisClient

    private usersInfoSuffix = 'users'
    private userChatIdsSuffix = 'chats'
    private userSocketIdSuffix = 'socket'
    private messagesSuffix = 'messages'

    private initConnection() {
        const host = 'redis'
        const port = (process.env.REDIS_PORT && process.env.REDIS_PORT.length > 0) 
            ? process.env.REDIS_PORT : '6379'
        const user = (process.env.REDIS_USER && process.env.REDIS_USER.length > 0) 
            ? process.env.REDIS_USER : 'chat_user'
        const password = (process.env.REDIS_PASSWORD && process.env.REDIS_PASSWORD.length > 0) 
            ? process.env.REDIS_PASSWORD : 'ffa9203c493aa99'

        const url = `redis://${user}:${password}@${host}:${port}`
        console.log("Redis url: ", url)

        this._client = redis.createClient({
            url: url
        })

        this._client.on("error", function(error) {
            console.error("Create client error: ", error);
        });
    }

    private addChatIdToUserChats(userId: string, chatId: string) {
        this._client.sadd(userId + ':' + this.userChatIdsSuffix, chatId)
    }

    private addChatRecord(userIds: string[]): string {
        const chatId = crypto
            .createHash('sha256')
            .update(userIds.concat(Date.now().toString()).join(' '))
            .digest('hex')

        this._client.sadd(chatId + ':' + this.usersInfoSuffix, userIds)
        return chatId
    }

    private addMessageIdToPool(messageId: string, chatId: string) {
        this._client.rpush(chatId + ':' + this.messagesSuffix, messageId)
    }

    private addMessageRecord(message: Message): string {
        const messageId = uuidv4()

        this._client.hmset(messageId,
            ['text', message.text,
                'state', message.state,
                'sender', message.sender,
                'recipient', message.recipient,
                'id', message.id,
                'date', message.date,
                'chatId', message.chatId])

        return messageId
    }

    private getUserChatIds(userId: string): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this._client.smembers(userId + ':' + this.userChatIdsSuffix, (err, res) => {
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
            this._client.smembers(chatId + ':' + this.usersInfoSuffix,
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
            this._client.hmget(chatId + ':' + this.messagesSuffix,
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
        this._client.set(userId + ':' + this.userSocketIdSuffix, socketId)
    }

    public unregisterUser(userId: string) {
        this._client.set(userId + ':' + this.userSocketIdSuffix, undefined)
    }

    public getUserSocketId(userId: string): Promise<string> {
        return new Promise<string>((resolve, reject) => {
            this._client.get(userId + ':' + this.userSocketIdSuffix, (err, res) => {
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

const DbClient = new ChatDbClient()
export default DbClient