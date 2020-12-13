import {Socket} from "socket.io";
const uuid = require("uuid").v4;


export class User {
    id: string;
    name: string;
    chats: string[];
    socket: Socket;

    constructor({
                    id,
                    name,
                    chats,
                    socket,
                }: {
        id: string;
        name: string;
        socket: Socket;
        chats?: string[];
    }) {
        this.id = id;
        this.name = name;
        this.socket = socket;
        this.chats = chats || [];
    }

    toString() {
        return JSON.stringify(this);
    }
}

export class WSmessage {
    type: string;
    payload: Chat | Message | null;

    constructor({type, payload}: { type: string; payload?: Chat | Message }) {
        this.type = type;
        this.payload = payload || null;
    }

    toString() {
        return JSON.stringify(this);
    }
}

export class Chat {
    id: string;
    userIds: string[];
    messages: Message[];

    constructor({
                    id,
                    userIds,
                    messages,
                }: {
        id?: string;
        userIds: string[];
        messages?: Message[];
    }) {
        this.id = id || uuid();
        this.userIds = userIds;
        this.messages = messages || [];
    }
}

export class Message {
    id: string;
    sender: string;
    recipient: string;
    date: number;
    text: string;
    status: string;
    chatId: string;

    constructor({
                    id,
                    sender,
                    recipient,
                    date,
                    text,
                    chatId,
                    status
                }: {
        id: string;
        sender: string;
        recipient: string;
        date: number;
        text: string;
        chatId: string;
        status: string;
    }) {
        this.id = id || uuid();
        this.chatId = chatId || "";
        this.sender = sender;
        this.recipient = recipient;
        this.date = date || new Date().getTime();
        this.text = text || "";
        this.status = status || "";
    }

    toString() {
        return JSON.stringify(this);
    }
}
