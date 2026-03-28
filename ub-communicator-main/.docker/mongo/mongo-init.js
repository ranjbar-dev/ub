
var db = connect("mongodb://localhost/admin");
db.createUser(
    {
        user: "ub_mongo_user",
        pwd: "ub_mongo_pass",
        roles: [
            {
                role: "readWrite",
                db: "ubMessages"
            }
        ]
    }
);