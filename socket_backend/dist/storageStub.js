"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getChats = exports.getUsers = void 0;
var chats = new Map();
var users = new Map();
function getUsers() {
    return users;
}
exports.getUsers = getUsers;
function getChats() {
    return chats;
}
exports.getChats = getChats;
