const MONGO_USER = process.env.MONGO_USER || "admin";
const MONGO_PASSWD = process.env.MONGO_PASSWD || "passwd";
const MONGO_ADDR = process.env.MONGO_ADDR || "localhost:27017";

export const DSN = `mongodb://${MONGO_USER}:${MONGO_PASSWD}@${MONGO_ADDR}/`
export const usersDb = `matcha`
export const usersCollection = `users`

console.log(DSN);