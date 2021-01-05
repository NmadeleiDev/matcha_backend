import mongo from 'mongodb';
import {DSN, usersCollection, usersDb} from './config'

export class MongoUser {
    constructor() {
        this.initConnection().catch(console.warn)
    }

    private _connection: mongo.MongoClient | undefined

    get connection(): mongo.MongoClient | undefined {
        return this._connection
    }

    set connection(val: mongo.MongoClient | undefined) {
        this._connection = val
    }

    async setUserOnlineState(id: string, state: boolean) {
        if (!this.connection)
            throw "Mongo client not connected yet!"

        try {
            const res = await this.connection
                .db(usersDb)
                .collection(usersCollection).updateOne({id: id}, {$set: {is_online: state}})
            if (res.modifiedCount !== 1) {
                console.log("FAILED TO SET ONLINE STATE FOR: ", id)
            }
        } catch (e) {
            console.log("Update error: ", e);
            throw(e)
        }
    }

    async updateUserLastOnlineTime(id: string) {
        if (!this.connection)
            throw "Mongo client not connected yet!"

        const time = Math.round(Date.now() / 1000)
        try {
            const res = await this.connection
                .db(usersDb)
                .collection(usersCollection).updateOne({id: id}, {$set: {last_online: time}})
            if (res.modifiedCount !== 1) {
                console.log("FAILED TO LAST ONLINE TIME FOR: ", id)
            }
        } catch (e) {
            console.log("Update error: ", e);
            throw(e)
        }
    }

    async initConnection() {
        try {
            const mongoClient = new mongo.MongoClient(DSN, {useUnifiedTopology: true});
            this._connection = await mongoClient.connect();
            console.log("Mongo connection succeded!");
        } catch (e) {
            console.log("Failed to connect to mongo: ", e);
            throw "Connection error!"
        }
    }
}

export const MongoManager = new MongoUser()
