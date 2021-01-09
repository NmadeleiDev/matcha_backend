const fs = require('fs');


async function deleteImagesFromStorage(filenames) {
    console.log("Deleting files: ", filenames)

    if (!Array.isArray(filenames)) {
        return null
    }

    for (let i = 0; i < filenames.length; i++) {
        const filename = filenames[i]
        if (typeof filename !== 'string')
            continue
        fs.unlink(filename, function(err){
            if(err) {
                console.log("Error deleting files: ", err)
            }
            console.log('Deleted! ', filenames);
        });
    }
}

exports.deleteImagesFromStorage = deleteImagesFromStorage;