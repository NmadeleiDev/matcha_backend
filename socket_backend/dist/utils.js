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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getAllUserChats = exports.mapUsersToIds = exports.mapChatsToIds = exports.getChatLength = exports.deleteUserFromChat = exports.getUserId = exports.deleteMessageFromChat = exports.updateMessageInChat = exports.findMessageInChat = exports.addMessageToChat = exports.sendToUser = void 0;
var model_1 = require("./model");
var config_1 = __importDefault(require("./config"));
var storage = __importStar(require("./storageStub"));
var sendToUser = function (WStype, type, payload, userIds) {
    if (WStype === config_1.default.WS.CHAT) {
        payload = payload;
    }
    else if (WStype === config_1.default.WS.MESSAGE) {
        payload = payload;
    }
    userIds.forEach(function (userId) {
        var user = storage.getUsers().get(userId);
        var wsMessage = new model_1.WSmessage({ type: type, payload: payload });
        console.log("[sendToUser]", { userId: userId, type: type, wsMessage: wsMessage });
        user && user.socket.emit(WStype, wsMessage.toString());
    });
};
exports.sendToUser = sendToUser;
var addMessageToChat = function (chatId, message) {
    var chat = storage.getChats().get(chatId);
    chat && chat.messages.push(message) && storage.getChats().set(chatId, chat);
};
exports.addMessageToChat = addMessageToChat;
var findMessageInChat = function (chatId, messageId) {
    var chat = storage.getChats().get(chatId);
    if (!chat) {
        console.log("Failed to find chat with id = " + chatId + " to find message = " + messageId);
        return null;
    }
    var message = chat.messages.find(function (item) { return item.id === messageId; });
    if (!message) {
        console.log("Failed to find message with id = " + messageId + " in chat = " + chatId);
        return null;
    }
    return message;
};
exports.findMessageInChat = findMessageInChat;
var updateMessageInChat = function (chatId, payload) {
    var message = exports.findMessageInChat(chatId, payload.id);
    if (!message)
        return;
    message.text = payload.text;
    message.status = payload.status;
    // chat.messages.push(payload) && storage.getChats().set(chatId, chat);
};
exports.updateMessageInChat = updateMessageInChat;
var deleteMessageFromChat = function (chatId, payload) {
    var chat = storage.getChats().get(chatId);
    if (!chat) {
        console.log("Failed to find chat with id = " + chatId + " to find message = " + payload.id);
        return null;
    }
    var idx = chat.messages.findIndex(function (item) { return item.id === payload.id; });
    if (idx < 0) {
        console.log("Failed to find message with id = " + payload.id);
        return;
    }
    chat.messages.splice(idx, 1);
};
exports.deleteMessageFromChat = deleteMessageFromChat;
var getUserId = function (socket) {
    var user = storage.getUsers().get(socket.id);
    var id = user && user.id;
    if (id)
        return id;
    else
        throw new Error("User not found by socket id");
};
exports.getUserId = getUserId;
var deleteUserFromChat = function (userId) {
    storage.getUsers().delete(userId);
};
exports.deleteUserFromChat = deleteUserFromChat;
var getChatLength = function (chatId) {
    var chat = storage.getChats().get(chatId);
    return chat ? chat.messages.length : 0;
};
exports.getChatLength = getChatLength;
var mapChatsToIds = function () {
    return storage.getChats() && storage.getChats().size ? Array.from(storage.getChats().keys()) : [];
};
exports.mapChatsToIds = mapChatsToIds;
var mapUsersToIds = function () {
    return storage.getUsers() && storage.getUsers().size ? Array.from(storage.getUsers().keys()) : [];
};
exports.mapUsersToIds = mapUsersToIds;
var getAllUserChats = function (userId) {
    console.log(__spread(storage.getChats().values()));
    __spread(storage.getChats().values()).forEach(function (el) { return console.log(el); });
    return userId && storage.getChats() && storage.getChats().size
        ? // if userId is in chat => return chat id
            __spread(storage.getChats().values()).filter(function (chat) { return !!chat.userIds.find(function (id) { return id === userId; }); })
                .map(function (chat) { return chat.id; })
        : [];
};
exports.getAllUserChats = getAllUserChats;
