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

    subscribe(userId: string) {
        if (!this.client) {
            console.error("redis client is undef for user: ", userId)
            return
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
}

export const notificationsClient = new NotificationsClient()
