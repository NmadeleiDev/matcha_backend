const pgp = require("pg-promise")();

function getPostgresDsn() {
    const user = process.env.POSTGRES_USER || "postgres";
    const passwd = process.env.POSTGRES_PASSWORD || "passwd";
    const host = process.env.POSTGRES_HOST || "localhost";
    const port = process.env.POSTGRES_PORT || "5432";
    const db = process.env.POSTGRES_DB || "postgres";

    return `postgres://${user}:${passwd}@${host}:${port}/${db}`
}

console.log(getPostgresDsn())
const db = pgp(getPostgresDsn());

async function getUserIdBySession(sessionId) {
    try {
        const data = await db.oneOrNone('SELECT id FROM user_data_schema.user_data WHERE session_key=$1', sessionId)
        return data.id
    } catch (e) {
        console.log("Error getting user id by session: ", e)
        return null
    }
}

exports.getUserIdBySession = getUserIdBySession;
