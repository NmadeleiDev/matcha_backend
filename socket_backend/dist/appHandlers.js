"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
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
exports.addServerHandlers = exports.addSocketHandlers = void 0;
var config_1 = __importDefault(require("./config"));
var model_1 = require("./model");
var utils = __importStar(require("./utils"));
var storage = __importStar(require("./storageStub"));
var utils_1 = require("./utils");
function addSocketHandlers(socket) {
    socket.on(config_1.default.WS.MESSAGE, function (message) {
        var json = JSON.parse(message);
        console.log(json);
        var payload = json.payload;
        switch (json.type) {
            case config_1.default.MESSAGE_TYPES.NEW_MESSAGE:
                // if first message in chat - send NEW_CHAT message first
                // save chat in db
                // add message to db
                // send to other NEW_MESSAGE
                // ---- or maybe just send DELIVERED_MESSAGE
                console.log("[NEW MESSAGE]");
                utils.addMessageToChat(payload.chatId, payload);
                payload.status = config_1.default.MESSAGE_STATUS.STATUS_DELIVERED;
                utils.sendToUser(config_1.default.WS.MESSAGE, config_1.default.MESSAGE_TYPES.NEW_MESSAGE, payload, [payload.recipient]);
                utils.sendToUser(config_1.default.WS.MESSAGE, config_1.default.MESSAGE_TYPES.UPDATE_MESSAGE, payload, [payload.sender]);
                return;
            case config_1.default.MESSAGE_TYPES.UPDATE_MESSAGE:
                utils.updateMessageInChat(payload.chatId, payload);
                utils.sendToUser(config_1.default.WS.MESSAGE, config_1.default.MESSAGE_TYPES.UPDATE_MESSAGE, payload, [payload.recipient]);
                return;
            case config_1.default.MESSAGE_TYPES.DELETE_MESSAGE:
                utils.deleteMessageFromChat(payload.chatId, payload);
                utils.sendToUser(config_1.default.WS.MESSAGE, config_1.default.MESSAGE_TYPES.UPDATE_MESSAGE, payload, [payload.recipient]);
                return;
            default:
                console.log("Unknown message type: ", json.type);
                return;
        }
    });
    socket.on(config_1.default.WS.CHAT, function (message) {
        var json = JSON.parse(message);
        var payload = json.payload;
        console.log("Chat type payload: ", json);
        switch (json.type) {
            case config_1.default.CHAT_TYPES.NEW_CHAT:
                var newChat_1 = new model_1.Chat(__assign({}, payload));
                // save chat in memory, wait for the first message
                console.log("[NEW_CHAT]");
                var chatExists = __spread(storage.getChats().values()).find(function (chat) {
                    return chat.userIds.includes(newChat_1.userIds[0]) &&
                        chat.userIds.includes(newChat_1.userIds[1]);
                });
                console.log(chatExists);
                if (chatExists) {
                    console.log("Chat to create already exists: ", chatExists);
                    utils_1.sendToUser(config_1.default.WS.CHAT, config_1.default.CHAT_TYPES.NEW_CHAT, chatExists, chatExists.userIds);
                    return;
                }
                storage.getChats().set(newChat_1.id, newChat_1);
                utils_1.sendToUser(config_1.default.WS.CHAT, config_1.default.CHAT_TYPES.NEW_CHAT, newChat_1, newChat_1.userIds);
                return;
            case config_1.default.CHAT_TYPES.DELETE_CHAT:
                // send to both DELETE_CHAT message
                // delete chat from db
                var chatToDelete = storage.getChats().get(payload.id);
                if (!chatToDelete) {
                    console.log("Error find chat to delete: ", chatToDelete);
                    return;
                }
                utils.sendToUser(config_1.default.WS.CHAT, config_1.default.CHAT_TYPES.DELETE_CHAT, chatToDelete, chatToDelete.userIds);
                storage.getChats().delete(payload.id);
                return;
            default:
                return;
        }
    });
    socket.on("error", function (e) {
        console.log(e);
    });
    socket.on("disconnect", function (reason) {
        console.log("disconnect: " + reason);
        try {
            utils.deleteUserFromChat(socket.id);
            console.log("user disconnected: " + utils.getUserId(socket));
        }
        catch (e) {
            console.log(e);
        }
    });
}
exports.addSocketHandlers = addSocketHandlers;
function addServerHandlers(app) {
    app.post("/chat", (function (req, res) {
        var newChat = new model_1.Chat({ userIds: req.body.userIds });
        var chatExists = __spread(storage.getChats().values()).find(function (chat) {
            return chat.userIds.includes(newChat.userIds[0]) &&
                chat.userIds.includes(newChat.userIds[1]);
        });
        console.log(chatExists);
        if (chatExists) {
            console.log("Chat to create already exists: ", chatExists);
            utils_1.sendToUser(config_1.default.WS.CHAT, config_1.default.CHAT_TYPES.NEW_CHAT, chatExists, chatExists.userIds);
            res.json(chatExists);
            return;
        }
        storage.getChats().set(newChat.id, newChat);
        utils_1.sendToUser(config_1.default.WS.CHAT, config_1.default.CHAT_TYPES.NEW_CHAT, newChat, newChat.userIds);
        res.json(newChat);
        return;
    }));
}
exports.addServerHandlers = addServerHandlers;
