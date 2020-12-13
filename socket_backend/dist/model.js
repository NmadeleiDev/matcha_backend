"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Message = exports.Chat = exports.WSmessage = exports.User = void 0;
var uuid = require("uuid").v4;
var User = /** @class */ (function () {
    function User(_a) {
        var id = _a.id, name = _a.name, chats = _a.chats, socket = _a.socket;
        this.id = id;
        this.name = name;
        this.socket = socket;
        this.chats = chats || [];
    }
    User.prototype.toString = function () {
        return JSON.stringify(this);
    };
    return User;
}());
exports.User = User;
var WSmessage = /** @class */ (function () {
    function WSmessage(_a) {
        var type = _a.type, payload = _a.payload;
        this.type = type;
        this.payload = payload || null;
    }
    WSmessage.prototype.toString = function () {
        return JSON.stringify(this);
    };
    return WSmessage;
}());
exports.WSmessage = WSmessage;
var Chat = /** @class */ (function () {
    function Chat(_a) {
        var id = _a.id, userIds = _a.userIds, messages = _a.messages;
        this.id = id || uuid();
        this.userIds = userIds;
        this.messages = messages || [];
    }
    return Chat;
}());
exports.Chat = Chat;
var Message = /** @class */ (function () {
    function Message(_a) {
        var id = _a.id, sender = _a.sender, recipient = _a.recipient, date = _a.date, text = _a.text, chatId = _a.chatId, status = _a.status;
        this.id = id || uuid();
        this.chatId = chatId || "";
        this.sender = sender;
        this.recipient = recipient;
        this.date = date || new Date().getTime();
        this.text = text || "";
        this.status = status || "";
    }
    Message.prototype.toString = function () {
        return JSON.stringify(this);
    };
    return Message;
}());
exports.Message = Message;
