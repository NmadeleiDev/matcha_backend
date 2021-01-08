import * as crypto from 'crypto'
import {Chat, Message} from "../../model";
import { v4 as uuidv4 } from 'uuid';
import {BaseRedis} from './base'
import {Socket} from "socket.io";
import * as utils from '../../utils'
import CONSTANTS from '../../config'

export class NotificationsClient extends BaseRedis{
    constructor() {
        super()

        this.operationalClient = new BaseRedis()

        if (!this.client) {
            console.error("redis client is undef!")
            return
        }

        this.client.on("subscribe", function(channel, count) {
            console.log(`User ${channel} subscribed to his channel: ${channel}, count=${count}`)
        });

        this.client.on("message", function(channel, message) {
            const body = JSON.parse(message) as {type: string; user: string}

            console.log(`User ${channel} got message from redis: ${message}`)
            utils.sendToUser(CONSTANTS.WS.UPDATE, body.type, {userId: body.user}, [channel])
        });
    }

    operationalClient: BaseRedis

    async subscribe(userId: string, doReadCache = true) {
        if (!this.client) {
            console.error("redis client is undef for user: ", userId)
            return
        }

        if (doReadCache) {
            while (true) {
                const message = await this.popValueFromArray(`${userId}:cached`)
                if (message === null)
                    break

                const body = JSON.parse(message) as {type: string; user: string}

                console.log(`User ${userId} got cached message in redis: ${message}`)
                utils.sendToUser(CONSTANTS.WS.UPDATE, body.type, {userId: body.user}, [userId])
            }
        }

        this.client.subscribe(userId)
    }

    unsubscribe(userId: string) {
        if (!this.client) {
            console.error("Client is undef !")
            return
        }

        this.client.unsubscribe(userId);
        // this.client.quit();
    }

    setOnlineState(user: string, state: boolean) {
        if (!this.operationalClient.isConnected) {
            console.error("Client is undef !")
            return
        }
        // @ts-ignore
        this.operationalClient.client.set(`${user}:online`, state ? '1' : '0')
    }

     popValueFromArray(key: string): Promise<string | null> {
        return new Promise((resolve, reject) => {
            if (!this.operationalClient.isConnected) {
                console.error("Client is undef !")
                resolve(null)
            }
            // @ts-ignore
            this.operationalClient.client.lpop(key, (err, val) => {
                if (err) {
                    console.log("Error pop val: ", err)
                    resolve(null)
                } else
                    resolve(val)
            })
        })
    }
}

export const notificationsClient = new NotificationsClient()
