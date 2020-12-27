import {Chat, User, WSmessage} from "./model";
import * as storage from './storageStub'

const queryString = require("query-string");
import express = require("express");
const bodyParser = require('body-parser').json()


const app = express();
import {Socket} from "socket.io";
import * as utils from "./utils";
import * as handlers from "./appHandlers";

const http = require("http").createServer(app);
app.use(bodyParser);
handlers.addServerHandlers(app);

const readline = require("readline");
// const jsonParser = require("socket.io-json-parser");

const port = process.env.BACKEND_PORT || '4444'

const io = require("socket.io")(http, {
    cors: {
        origin: `http://localhost:3000`,
        credentials: true,
    },
    // transports: ['websocket'],
    path: "/connect",
    // parser: jsonParser,
});

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
});


// console.log(io);

// ===========================
// ======== SOCKET ===========
// ===========================

io.on("connection", (socket: Socket) => {
    // console.log("New conn: ", socket);

    const parsed = queryString.parse(socket.handshake.url.split("?")[1]);
    console.log("new conn url parsed:", parsed);
    const userId = parsed.auth;
    const user = new User({
        id: userId,
        name: userId,
        socket: socket,
    });
    storage.getUsers().set(userId, user)

    handlers.setOnlineState(userId, true)
    handlers.addSocketHandlers(userId, socket)
});

http.listen(parseInt(port), () => {
    console.log(`listening on *:${port}`);
});

rl.on("line", (input: string) => {
    if (input.startsWith("users")) {
        console.log(storage.getUsers());
    } else if (input.startsWith("user ids")) {
        console.log(utils.mapUsersToIds());
    } else if (input.startsWith("chats")) {
        console.log(storage.getUsers());
    } else if (input.startsWith("chat ids")) {
        console.log(utils.mapChatsToIds());
    }
});

