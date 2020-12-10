"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const http_1 = __importDefault(require("http"));
const socket_io_1 = __importDefault(require("socket.io"));
const express_1 = __importDefault(require("express"));
const main_1 = require("./handlers/main");
const body_parser_1 = __importDefault(require("body-parser"));
const app = express_1.default();
app.use(body_parser_1.default.json());
const server = http_1.default.createServer(app);
const io = new socket_io_1.default.Server(server, {
    path: '/ws',
    serveClient: false,
    pingInterval: 10000,
    pingTimeout: 5000,
    cookie: false,
    cors: { origin: true }
});
app.get('/', (req, res) => {
    res.send('<h1>Start chatting!</h1>');
});
app.post('/chat', (req, res) => {
    const chat = req.body;
    main_1.sendChatCreationMessage(chat);
    res.end("Success");
});
io.on('connection', (socket) => {
    socket.auth;
    socket.on('chat message', (msg) => {
        console.log('message: ' + msg);
    });
});
const port = (process.env.BACKEND_PORT && process.env.BACKEND_PORT.length > 0)
    ? parseInt(process.env.BACKEND_PORT) : 4444;
server.listen(port);
//# sourceMappingURL=app.js.map