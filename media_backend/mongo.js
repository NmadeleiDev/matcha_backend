const mongo = require('mongodb');

const MONGO_USER = process.env.MONGO_USER || "admin";
const MONGO_PASSWORD = process.env.MONGO_PASSWORD || "passwd";
const MONGO_ADDRESS = process.env.MONGO_ADDRESS || "localhost:27017";

const dsn = `mongodb://${MONGO_USER}:${MONGO_PASSWORD}@${MONGO_ADDRESS}/`

var client

console.log(dsn);

function mongoFind(collection, findObj) {
    return new Promise((resolve, reject) => {
        collection.find(findObj)
            .toArray((err, items) => {
                if (!Array.isArray(items) || err !== null) {
                    console.log("find error: ", err)
                    reject(err)
                    return
                }
                resolve(items)
            })
    })
}


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
    // return wsClient;
}

async function closeConnection() {
    try {
        await client.close()
    } catch (e) {
        console.log("Failed to close mongo: ", e);
    }
}

async function insertImageDataToMediaCollection(data) {
    const mediaCollection = client.db("media").collection("images");
    try {
        const result = await mediaCollection.insertOne(data)
        return result.insertedId.toString()
    } catch (e) {
        console.log("error inserting image data: ", e)
        return null
    }
}

async function getUserImageIds(userId) {
    const userCollection = client.db("matcha").collection("users");

    try { // TODO: projection не срабатывает!
        let result = await userCollection.findOne({id: userId}, {images: 1}).images
        if (Array.isArray(result))
            return result
        else
            return []
    } catch (e) {
        console.log(e)
        return []
    }
}

async function insertImageData(data) {
    let update;

    const userCollection = client.db("matcha").collection("users");

    if ((await getUserImageIds(data.id)).length >= 5) {
        return false
    }
    console.log("Passed num images limit")

    let insertedId = await insertImageDataToMediaCollection(data);
    if (!insertedId) {
        return null
    }
    if (typeof insertedId !== 'string') {
        console.log("Inserted id is not string!", typeof insertedId, insertedId, data)
        return null
    }
    update = {$push: {images: insertedId}}

    if (data.isAvatar === true || data.isAvatar === 'true')
        update.$set = {avatar: insertedId}

    try {
        await userCollection.updateOne({id: data.id}, update)
    } catch (e) {
        console.log("error inserting image data: ",e)
        return null
    }
    return insertedId;
}

async function setUserAvatar(userId, imageId) {
    try {
        const result = await getUserImageIds(userId)
        console.log("Got user images: ", result)
        if (!result.some(item => item === imageId)) {
            return false
        }
    } catch (e) {
        console.log(e)
        return null
    }
    console.log(`Image ${imageId} for ${userId} found`)

    const userCollection = client.db("matcha").collection("users");

    let update = {$set: {avatar: imageId}}
    try {
        await userCollection.updateOne({id: userId}, update)
    } catch (e) {
        console.log(e)
        return null
    }
    return true
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

async function deleteImagesIdFromUserImages(imageIds, userId) {
    const update = {$pullAll: {images: imageIds}}

    const userCollection = client.db("matcha").collection("users")

    console.log("Updating: ", userId, update)
    try {
        const res = await userCollection.updateOne({id: userId}, update)
        console.log("Update (delete images) res: ", res.result.nModified)
        return res.result.nModified === imageIds.length
    } catch (e) {
        console.log(e)
        return false
    }
}

async function deleteImageData(imageIds) {
    const collection = client.db("media").collection("images");
    const images = []
    const filter = { _id: { $in: imageIds.map(id => new mongo.ObjectID(id)) }}

    try {
        images.push(...(await collection.find(filter).toArray()))
    } catch (e) {
        console.log("Error finding images: ", e)
        return null
    }
    if (images.length < 1)
        return null
    const userId = images[0].id

    if (!(await deleteImagesIdFromUserImages(images.map(item => item._id.toString()), userId))) {
        return null
    }

    try {
        await collection.deleteMany(filter)
    } catch (e) {
        console.log("Error deleting images: ", e)
        return null
    }

    return images.map(item => item.filename)
}

exports.initConnection = initConnection;
exports.closeConnection = closeConnection;
exports.getFileByDocumentId = getFileByDocumentId;
exports.insertImageData = insertImageData;
exports.deleteImageData = deleteImageData;
exports.setUserAvatar = setUserAvatar;