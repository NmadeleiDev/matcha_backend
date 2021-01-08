import {Socket} from "socket.io";
import {User} from "./model";
import * as storage from "./storageStub";
import * as handlers from "./appHandlers";
import * as utils from './utils'
import {notificationsClient} from "./db/redis/notifications";

export function initClient(userId: string, socket: Socket) {
    const user = new User({
        id: userId,
        name: userId,
        socket: socket,
    });
    storage.getUsers().set(userId, user)
    utils.setOnlineState(userId, true)

    notificationsClient.subscribe(userId).catch((e) => console.warn(e))
    handlers.addSocketHandlers(userId, socket)
}