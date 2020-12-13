import {Chat, User} from "./model";

const chats = new Map<string, Chat>();
const users = new Map<string, User>();

export function getUsers(): Map<string, User> {
    return users
}

export function getChats(): Map<string, Chat> {
    return chats
}