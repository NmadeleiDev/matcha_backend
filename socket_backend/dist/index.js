"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
var __spread = (this && this.__spread) || function () {
    for (var ar = [], i = 0; i < arguments.length; i++) ar = ar.concat(__read(arguments[i]));
    return ar;
};
Object.defineProperty(exports, "__esModule", { value: true });
var model_1 = require("./model");
var storage = __importStar(require("./storageStub"));
var queryString = require("query-string");
var express = require("express");
var bodyParser = require('body-parser');
var app = express();
var utils = __importStar(require("./utils"));
var handlers = __importStar(require("./appHandlers"));
var http = require("http").createServer(app);
app.use(bodyParser);
handlers.addServerHandlers(app);
var readline = require("readline");
// const jsonParser = require("socket.io-json-parser");
var port = process.env.BACKEND_PORT || '4444';
var io = require("socket.io")(http, {
    cors: {
        origin: "http://localhost:3000",
        credentials: true,
    },
    path: "/ws",
});
var rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
});
// console.log(io);
// ===========================
// ======== SOCKET ===========
// ===========================
io.on("connection", function (socket) {
    // console.log(socket);
    var parsed = queryString.parse(socket.handshake.url.split("?")[1]);
    console.log(parsed);
    var userId = parsed.auth;
    var user = new model_1.User({
        id: userId,
        name: userId,
        socket: socket,
    });
    storage.getUsers().set(userId, user);
    // debug only
    socket.onAny(function (event) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        return console.log.apply(console, __spread([event], args));
    });
    handlers.addSocketHandlers(socket);
});
http.listen(parseInt(port), function () {
    console.log("listening on *:" + port);
});
rl.on("line", function (input) {
    if (input.startsWith("users")) {
        console.log(storage.getUsers());
    }
    else if (input.startsWith("user ids")) {
        console.log(utils.mapUsersToIds());
    }
    else if (input.startsWith("chats")) {
        console.log(storage.getUsers());
    }
    else if (input.startsWith("chat ids")) {
        console.log(utils.mapChatsToIds());
    }
});
