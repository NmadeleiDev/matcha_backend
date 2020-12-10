export interface Chat {
    id: string
    userIds: string[]
    messages: Message[]
}

export interface Message {
    id: string
    text: string
    sender: string
    recipient: string
    date: number
    state: number
    chatId: string
}