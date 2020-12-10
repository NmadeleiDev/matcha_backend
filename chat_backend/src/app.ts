import http from 'http';
import socket from 'socket.io';
import express from 'express';
import {registerSocketConn, sendChatCreationMessage} from "./socketUtils/socket";
import {Chat} from "./model/model";

import bodyParser from 'body-parser';
import { setHandlers } from './handlers/main';


const app = express()
app.use(bodyParser.json())

const server = http.createServer(app)
const io = new socket.Server(server, {
    path: '/connect',
    serveClient: false,
    pingInterval: 10000,
    pingTimeout: 5000,
    cookie: false,
    cors: { origin: true }
});

setHandlers(app)

io.on('connection', (socket) => {
    console.log("Auth: ", socket.auth);
    
    registerSocketConn(socket, socket.auth.id)
    // socket.on('chat message', (msg) => {
    //     console.log('message: ' + msg);
    // });
});

const port = (process.env.BACKEND_PORT && process.env.BACKEND_PORT.length > 0) 
    ? parseInt(process.env.BACKEND_PORT) : 4444

server.listen(port);