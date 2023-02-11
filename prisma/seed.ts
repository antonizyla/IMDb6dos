import { randomUUID } from "crypto"

import axios from 'axios'
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');

const { PrismaClient } = require('@prisma/client')
const prisma = new PrismaClient()


async function fetchAndUnzip(filename: string, label = undefined) {

    if (fs.existsSync(path.join(__dirname, `data/${filename}.tsv`))) {
        console.log(`${label || filename} File Exists`)
    } else {
        console.log(`Fetching ${label || filename} file`)
        const tsvData = await axios.get(`https://datasets.imdbws.com/${filename}.tsv.gz`, { responseType: 'arraybuffer', 'decompress': true }).then(res => {
            return res.data;
        }).catch((err: any) => {
            console.log(err);
            return;
        });
        console.log(`Fetched ${label || filename} file`)
        fs.writeFileSync(path.join(__dirname, `data/${filename}.tsv.gz`), tsvData);

        // unzip file
        const contents = fs.createReadStream(path.join(__dirname, `data/${filename}.tsv.gz`));
        const writestream = fs.createWriteStream(path.join(__dirname, `data/${filename}.tsv`));
        const unzip = zlib.createGunzip();
        let stream = contents.pipe(unzip).pipe(writestream);
        console.log(`Unzipping ${label || filename} file`)

        stream.on('finish', () => {
            console.log(`Finished Unzipping ${label || filename} file`)
            fs.unlinkSync(path.join(__dirname, `data/${filename}.tsv.gz`));
        })
    }
}


async function main() {
    const files = ['title.ratings', 'title.akas', 'title.basics', 'title.episode', 'title.principals', 'name.basics'];
    files.forEach(async (file) => {
        fetchAndUnzip(file);
    });
}

main().then(async () => {
    await prisma.$disconnect()
}).catch(async e => {
    console.error(e)
    await prisma.$disconnect()
    process.exit(1)
});
