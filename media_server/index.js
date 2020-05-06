const express = require('express');
const multer  = require('multer')
const bodyParser = require('body-parser');
const mongoFuncs = require('./mongo');

const MEDIA_SERVER_PORT = process.env.MEDIA_SERVER_PORT || 3333;
const STORAGE_PATH = process.env.STORAGE_PATH || '/app/storage/';

const app = express();
app.use(bodyParser.json());

mongoFuncs.initConnection().catch(console.log);

const Storage = multer.diskStorage({
    destination: function(req, file, callback) {
        callback(null, STORAGE_PATH);
    },
    filename: function(req, file, callback) {
        callback(null, Date.now() + "_" + Math.floor(Math.random() * 1000) + '_' + file.originalname.replace(/ /g, ''));
    }
});

const upload = multer({ dest: 'uploads/', storage: Storage })

app.post("/upload", upload.single('image'), function (req, res, next) {
    // req.file is the `avatar` file
    console.log(req.file);
    const fileData = {
        filename: req.file.filename,
        email: req.body.email,
    }
    mongoFuncs.insertImageData(fileData).then(id => {
        if (id === null) {
            console.log("Error saving file data");
            res.end(JSON.stringify({status: false, data: "Error saving file data"}));
        } else {
            res.end(JSON.stringify({status: true, data: {id: id}}));
        }
    }).catch(console.log);
});

app.get("/img/*", function (req, res) {
    const fileId = req.params[0];
    mongoFuncs.getFileByDocumentId(fileId).then(data => {
        console.log(data);
        res.sendFile( STORAGE_PATH + data.filename);
    }).catch(console.log);
});

app.listen(MEDIA_SERVER_PORT, function(a) {
    console.log("Listening on port " + MEDIA_SERVER_PORT);
});

// mongoFuncs.closeConnection(mongoClient).catch(console.log);