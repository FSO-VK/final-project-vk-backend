const dbName = process.env.MEDICATION_MONGO_DATABASE;
const username = process.env.MEDICATION_MONGO_USER;
const password = process.env.MEDICATION_MONGO_PASSWORD;

if (!dbName || !username || !password) {
    console.error('Environment variables MONGODB_DATABASE, MONGODB_USER and MONGODB_PASSWORD must be set');
    process.exit(1);
}

db = db.getSiblingDB(dbName);
db.createUser({
    'user': username,
    'pwd': password,
    'roles': [
        {
            'role': 'readWrite',
            'db': dbName
        }
    ]
});

console.log('MongoDB user created');