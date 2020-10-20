const express = require('express');
const multer  = require('multer')
const cookieParser = require("cookie-parser");
const bodyParser = require('body-parser');
// const cors = require('cors');
const mongoFuncs = require('./mongo');
const pgFuncs = require('./postgres');

const MEDIA_SERVER_PORT = process.env.MEDIA_SERVER_PORT || 3333;
const STORAGE_PATH = process.env.STORAGE_PATH || '/app/storage/';

const app = express();
app.use(bodyParser.json());
app.use(cookieParser())
// app.use(cors());

mongoFuncs.initConnection().catch(console.log);

const Storage = multer.diskStorage({
    destination: function(req, file, callback) {
        callback(null, STORAGE_PATH);
    },
    filename: function(req, file, callback) {
        callback(null, Date.now() + "_" + Math.floor(Math.random() * 1000) + '_' + file.fieldname.replace(/ /g, ''));
    }
});

const upload = multer({ dest: 'uploads/', storage: Storage })

app.post("/upload", upload.single('user_image'), function (req, res, next) {
    if (!req.file) {
        res.end("File is empty");
        return;
    }

    const sessionId = req.cookies.session_id
    if (!sessionId) {
        res.end("Session cookie is empty");
        return;
    }

    pgFuncs.getUserIdBySession(sessionId).then(async (userId) => {
        if (!userId) {
            res.end(JSON.stringify({status: false, data: `Session key ${sessionId} is not valid`}));
            return;
        }
        const fileData = {
            filename: req.file.filename,
            id: userId,
            isAvatar: !(req.body.is_avatar) ? false : req.body.is_avatar,
        }
        try {
            const id = await mongoFuncs.insertImageData(fileData)
            res.end(JSON.stringify({status: true, data: {id: id}}));
        } catch (e) {
            res.end(JSON.stringify({status: false, data: "Error saving file data"}));
        }
    })
});

app.get("/img/*", function (req, res) {
    const fileId = req.params[0];
    mongoFuncs.getFileByDocumentId(fileId).then(data => {
        console.log(data);
        res.sendFile(STORAGE_PATH + data.filename);
    }).catch(console.log);
});

app.listen(MEDIA_SERVER_PORT, function(a) {
    console.log("Listening on port " + MEDIA_SERVER_PORT);
});

// mongoFuncs.closeConnection(mongoClient).catch(console.log);