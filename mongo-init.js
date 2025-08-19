// mongo-init.js
// This script will create the 'chitauri' database and insert 10 random API keys into the 'api_keys' collection.


const keys = [
    { key: "d3f8a1c2e4b5f6a7d8c9e0b1a2f3c4d5" },
    { key: "a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4" },
    { key: "f1e2d3c4b5a6f7e8d9c0b1a2e3f4d5c6" },
    { key: "c8d7e6f5a4b3c2d1e0f9a8b7c6d5e4f3" },
    { key: "b2a3c4d5e6f7a8b9c0d1e2f3a4b5c6d7" },
    { key: "e5f4d3c2b1a0e9f8d7c6b5a4e3f2d1c0" },
    { key: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6" },
    { key: "f9e8d7c6b5a4f3e2d1c0b9a8e7f6d5c4" },
    { key: "c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8" },
    { key: "b6a5c4d3e2f1b0a9c8d7e6f5a4b3c2d1" }
];

print('Inserting API keys:', keys);
db = db.getSiblingDB('chitauri');
db.api_keys.insertMany(keys);
