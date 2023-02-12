import { randomUUID } from "crypto"
import axios from 'axios'
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');
const es = require('event-stream');

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

async function seed_titles() {
    // read file using stream to avoid memory issues
    let thousandQueries: any[] = [];
    let s = fs.createReadStream(__dirname + `/data/title.basics.tsv`).pipe(es.split()).pipe(es.mapSync(async function(line: string) {
        let splitLine = line.split('\t');
        if (splitLine[0] != 'tconst') {

            const data = {
                title_id: splitLine[0],
                title_type: splitLine[1],
                primary_title: splitLine[2],
                original_title: splitLine[3],
                is_adult: Boolean(parseInt(splitLine[4])),
                start_year: parseInt(splitLine[5]) || 0,
                end_year: parseInt(splitLine[6]) || 0,
                runtime_minutes: parseInt(splitLine[7]) || 0,
                genres: splitLine[8].split(','),
            }
            thousandQueries.push(data);
            if (thousandQueries.length === 10000) {
                await prisma.titles.createMany({
                    data: thousandQueries,
                    skipDuplicates: true,
                });
                console.log("inserted 10000");
                thousandQueries = new Array();
            }
        }
    }));
    await prisma.titles.createMany({
        data: thousandQueries,
        skipDuplicates: true,
    });
}


async function main() {
    const files = ['title.ratings', 'title.akas', 'title.basics', 'title.episode', 'title.principals', 'name.basics'];
    files.forEach(async (file) => {
        fetchAndUnzip(file);
    });
    await seed_titles();
    //
    /*const title = await prisma.titles.create({
        data: {
            title_id: randomUUID(),
            title_type: 'movie',
            primary_title: 'test',
            original_title: 'test',
            is_adult: false,
            start_year: 2021,
            end_year: 0,
            runtime_minutes: 0,
        }
    });
    */
}

main().then(async () => {
    await prisma.$disconnect()
}).catch(async e => {
    console.error(e)
    await prisma.$disconnect()
    process.exit(1)
});
