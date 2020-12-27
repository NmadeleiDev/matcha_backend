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
    }

    subscribe(userId: string, socket: Socket) {
        if (!this.client) {
            console.error("Client is undef for user: ", userId)
            return
        }

        this.client.on("subscribe", function(channel, count) {
            console.log(`User ${userId} subscribed to his channel: ${channel}`)
        });

        this.client.on("message", function(channel, message) {
            console.log(`User ${userId} got message from redis: ${message.toString()}`)
            utils.sendToUser(CONSTANTS.WS.UPDATE, CONSTANTS.UPDATE_TYPES.NEW_LIKE, {data: message}, [userId])
        });
        this.client.subscribe(userId)
    }

    unsubscribe() {
        if (!this.client) {
            console.error("Client is undef !")
            return
        }

        this.client.unsubscribe();
        this.client.quit();
    }
}

export const notificationsClient = new NotificationsClient()
