const fs = require('fs');


async function deleteImagesFromStorage(filenames) {
    console.log("Deleting files: ", filenames)

    for (const filename in filenames) {
        if (typeof filename !== 'string')
            continue
        try {
            fs.unlinkSync(filename)
        } catch (e) {
            console.log(e)
        }
    }
}

exports.deleteImagesFromStorage = deleteImagesFromStorage;