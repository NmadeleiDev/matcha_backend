"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.sendChatCreationMessage = exports.registerSocketConn = void 0;
const main_1 = __importDefault(require("../db/main"));
const Clients = new Map();
function registerSocketConn(socket, clientId) {
    main_1.default.registerUserAsOnline(clientId, socket.id);
    Clients.set(clientId, socket);
    main_1.default.getUserChats(clientId)
        .then(chats => {
        socket.send(chats);
        chats.forEach(chat => {
            socket.join(chat.id); // подключаем пользователя к его чатам
        });
    });
    socket.on('message', (msg) => {
        console.log('message: ', msg);
        const message = msg;
        message.state = 2;
        main_1.default.addMessageToChat(message);
        socket.to(message.chatId).emit('message', message); // под каждый чат создается комната, к которой подключены оналайн пользователи
    });
}
exports.registerSocketConn = registerSocketConn;
function sendChatCreationMessage(chat) {
    chat.userIds.forEach(id => {
        const client = Clients.get(id);
        if (client) {
            client.emit('chat', chat);
        }
        else {
            console.log("Client is not online: ", id);
        }
    });
}
exports.sendChatCreationMessage = sendChatCreationMessage;
//# sourceMappingURL=main.js.map