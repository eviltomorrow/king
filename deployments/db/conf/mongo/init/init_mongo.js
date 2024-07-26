db = db.getSiblingDB('transaction_db');
db.createUser({"user":"admin","pwd":"admin123","roles":[{"role":"dbOwner","db":"transaction_db"}]});
db.createCollection('metadata');
db.metadata.createIndex({date: 1, code: 1},{background: true});