const fs = require('fs');


async function deleteImagesFromStorage(filenames) {
    for (const filename in filenames) {
        try {
            fs.unlinkSync(filename)
        } catch (e) {
            console.log(e)
        }
    }
}

exports.deleteImagesFromStorage = deleteImagesFromStorage;