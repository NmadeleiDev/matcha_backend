const express = require('express');
const multer  = require('multer')
const cookieParser = require("cookie-parser");
const bodyParser = require('body-parser');
// const cors = require('cors');
const mongoFuncs = require('./mongo');
const pgFuncs = require('./postgres');
const utils = require('./utils');

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

const upload = multer({ storage: Storage })

app.post("/upload", upload.single('userImage'), function (req, res) {
    if (!req.file) {
        res.end({status: false, data: "File is empty"});
        return;
    }

    const sessionId = req.cookies.session_id
    if (!sessionId) {
        res.end({status: false, data: "Session cookie is empty"});
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
            isAvatar: !(req.body.isAvatar) ? false : req.body.isAvatar,
        }
        try {
            const result = await mongoFuncs.insertImageData(fileData)
            if (result === false) {
                res.end(JSON.stringify({status: false, data: 'images limit reached'}));
            } else if (!result) {
                res.end(JSON.stringify({status: false, data: 'error saving image'}));
            } else {
                res.end(JSON.stringify({status: true, data: {id: result}}));
            }
        } catch (e) {
            res.end(JSON.stringify({status: false, data: "Error saving file data"}));
        }
    })
});

app.put("/avatar", function (req, res) {
    const sessionId = req.cookies.session_id
    if (!sessionId) {
        res.end({status: false, data: "Session cookie is empty"});
        return;
    }
    const imageId = req.body.imageId
    if (!imageId) {
        res.end({status: false, data: "imageId field is not valid"});
        return;
    }
    pgFuncs.getUserIdBySession(sessionId).then(async (userId) => {
        if (!userId) {
            res.end(JSON.stringify({status: false, data: `Session key ${sessionId} is not valid`}));
            return
        }
        let setResult = await mongoFuncs.setUserAvatar(userId, imageId)
        if (setResult === true) {
            res.end(JSON.stringify({status: true, data: "Avatar set successfully"}));
        } else if (setResult === false) {
            res.end(JSON.stringify({status: false, data: "imageId is not in user images"}));
        } else {
            res.end(JSON.stringify({status: false, data: "error setting user avatar"}));
        }
    })
})

app.delete("/img/:id", function (req, res) {
    const response = {
        status: false,
        data: 'error'
    }
    const imageId = req.params.id

    mongoFuncs.deleteImageData([imageId]).then(result => {
        console.log("Delete result: ", result)
        if (Array.isArray(result)) {
            response.data = 'success'
            response.status = true
            res.end(JSON.stringify(response))
            utils.deleteImagesFromStorage(result.map(item => STORAGE_PATH + item)).catch(console.error)
        } else {
            res.end(JSON.stringify(response))
        }
    }).catch((e) => console.error(e))
})

app.get("/img/:id", function (req, res) {
    const fileId = req.params.id;
    mongoFuncs.getFileByDocumentId(fileId).then(data => {
        if (!data) {
            res.json({status: false, data: 'File not found'}).status(404)
            return
        }
        console.log("Got image data: ", data);
        res.sendFile(STORAGE_PATH + data.filename);
    }).catch("Error gettign image data: ", console.log);
});

app.listen(MEDIA_SERVER_PORT, function(a) {
    console.log("Listening on port " + MEDIA_SERVER_PORT);
});

// mongoFuncs.closeConnection(mongoClient).catch(console.log);