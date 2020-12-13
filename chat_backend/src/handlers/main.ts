import { Application } from "express";
import { sendChatCreationMessage } from "../socketUtils/socket";
import { Chat } from "../model/model";
import DbClient from '../db/redis'

export function setHandlers(app: Application) {
    app.get('/', (req, res) => {
        res.send('<h1>Start chatting!</h1>');
    });
    
    app.post('/chat', (req, res) => {
        const chat = req.body as Chat

        console.log("Got chat create body: ", chat);
    
        DbClient.createChat(chat.userIds)
        sendChatCreationMessage(chat)
    
        res.end("Success")
    })
}