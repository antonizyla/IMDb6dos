import { randomUUID } from "crypto"

const { PrismaClient } = require('@prisma/client')

const prisma = new PrismaClient()

const user = {
    name: 'Alice',
    email: `${randomUUID()}@gmail.com`
}

async function main() {
    console.log(`Start seeding ...`)
    const createdUser = await prisma.user.create({
        data: user,
    })
    console.log(`Created user with id: ${createdUser.id}`)
}

main().then(async () => {
    await prisma.$disconnect()
}).catch(async e => {
    console.error(e)
    await prisma.$disconnect()
    process.exit(1)
});
