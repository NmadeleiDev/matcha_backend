const mongo = require('mongodb');

const MONGO_USER = process.env.MONGO_USER || "admin";
const MONGO_PASSWORD = process.env.MONGO_PASSWORD || "passwd";
const MONGO_ADDRESS = process.env.MONGO_ADDRESS || "localhost:27017";

const dsn = `mongodb://${MONGO_USER}:${MONGO_PASSWORD}@${MONGO_ADDRESS}/`

var client

console.log(dsn);

async function initConnection() {
    try {
        const mongoClient = new mongo.MongoClient(dsn, { useUnifiedTopology: true, useNewUrlParser: true });
        client = await mongoClient.connect();
    } catch (e) {
        console.log("Failed to connect to mongo: ", e);
        return null;
    } finally {
        console.log("All good");
    }
    // return client;
}

async function closeConnection() {
    try {
        await client.close()
    } catch (e) {
        console.log("Failed to close mongo: ", e);
    }
}

async function insertImageData(data) {
    let result;

    const mediaCollection = client.db("media").collection("images");
    try {
        result = await mediaCollection.insertOne(data)
    } catch (e) {
        console.log("error inserting image data: ",e)
        return null;
    }
    const userCollection = client.db("matcha").collection("users");
    try {
        await userCollection.updateOne({email: data.email}, {$push: {images: result.insertedId.toString()}})
    } catch (e) {
        console.log("error inserting image data: ",e)
        return null;
    }
    return result.insertedId;
}

async function getFileByDocumentId(id) {
    let result;

    const collection = client.db("media").collection("images");
    try {
        result = await collection.findOne({'_id': new mongo.ObjectID(id)});
    } catch (e) {
        console.log("error inserting image data: ",e)
        return null;
    }
    return result;
}

exports.initConnection = initConnection;
exports.closeConnection = closeConnection;
exports.getFileByDocumentId = getFileByDocumentId;
exports.insertImageData = insertImageData;