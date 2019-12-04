//put interpreter into strict mode
"use strict";

//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");
const mongo = require("mongodb");
const ObjectId = require("mongodb").ObjectId;
var MongoClient = require("mongodb").MongoClient;
var url = "mongodb://mongodb:27017/";
const contentType = "Content-Type";
const appJson = "application/json";
const amqp = require('amqplib/callback_api')

MongoClient.connect(url, function (err, db) {
  if (err) throw err;
  let dbo = db.db("mydb");
  console.log("Database created!");
  dbo.createCollection("listings", (err, res) => {
    console.log("Listing collection created");
  })
});

//create a new express application
const app = express();

//get ADDR environment variable,
//defaulting to ":80"
const addr = process.env.NODE_ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

let channel;
amqp.connect("amqp://" + process.env.RABBITMQADDR + ":5672/", (err, conn) => {
  conn.createChannel(function (err, ch) {
    var q = process.env.RABBITMQADDR;

    ch.assertQueue(q, { durable: false });
    channel = ch;
  });
});

app.route("/v1/listings")
  .get((req, res) => {
    MongoClient.connect(url, (err, db) => {
      if (err) throw err;
      let dbo = db.db("mydb");
      dbo
        .collection("listings")
        .find({})
        .toArray(function (err, result) {
          if (err) throw err;
          res.status(200);
          res.set(contentType, appJson);
          res.json(result);
          db.close();
        });
    });
  })
app.route("/v1/listings").post((req, res) => {
  if (!req.get("X-User") || req.get("X-User").length == 0) {
    res.status(401).send("Unauthorized");
    return;
  }
  let reqUserID = JSON.parse(req.header("X-User")).id;
  MongoClient.connect(url, function (err, db) {
    if (err) throw err;
    let dbo = db.db("mydb");
    let listingObj = req.body;
    listingObj.title = req.body.title;
    listingObj.description = req.body.description;
    listingObj.condition = req.body.condition;
    listingObj.price = req.body.price;
    listingObj.location = req.body.location;
    listingObj.contact = req.body.contact;
    listingObj.creator = reqUserID;
    dbo.collection("listings").insertOne(listingObj, async function (err, result) {
      if (err) throw err;
      db.close();
      const event = {
        type: "listing-new",
        message: listingObj
      };
      channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
      res.status(201);
      res.set("Content-Type", "application/json");
      res.json(listingObj);
    });
  });
})
app.route("/v1/listings/:id")
  .get((req, res) => {
    // if (!req.get("X-User") || req.get("X-User").length == 0) {
    //   res.status(401).send("Unauthorized");
    //   return;
    // }
    if (req.params.id.length != 24) {
      res.status(403).send("Invalid listing ID");
      return;
    }
    MongoClient.connect(url, (err, db) => {
      if (err) throw err;
      let dbo = db.db("mydb");
      dbo
        .collection("listings")
        .findOne({ _id: new ObjectId(req.params.id) }, (err, result) => {
          if (err) throw err;
          if (!result) {
            res.sendStatus(404);
            db.close();
            return;
          }
          res.status(200);
          res.set(contentType, appJson);
          res.json(result);
          db.close();
        })
    })
  })
app.route("/v1/listings/:id")
.patch((req, res) => {
  if (!req.get("X-User") || req.get("X-User").length == 0) {
    res.status(401).send("Unauthorized");
    return
  }
  if (req.params.id.length != 24) {
    res.status(403).send("Invalid listing ID");
    return
  }
  let listingID = new mongo.ObjectId(req.params.id);
  let reqUserID = JSON.parse(req.header("X-User")).id;
  let reqBody = req.body;
  MongoClient.connect(url, async (err, db) => {
    if (err) throw err;
    let dbo = db.db("mydb");
    const result = await dbo.collection("listings").findOne({ _id: listingID });
    if (!result || reqUserID != result.creator) {
      res.status(403).send("Unauthorized user or channel does not exist");
    } else {
      dbo
        .collection("listings")
        .findOneAndUpdate(
          { _id: listingID },
          { $set: { title: reqBody.title, description: reqBody.description, condition: reqBody.condition, location: reqBody.location, contact: reqBody.contact, price: reqBody.price } },
          { returnOriginal: false },
          (err, document) => {
            if (err) throw err;
            let result = document.value;
            let channelMembers = [];
            result.members && result.private && result.members.forEach(member => {
              channelMembers.push(member.id);
            });
            const event = {
              type: "channel-update",
              channel: result,
              userIDs: channelMembers
            };
            // channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
            res.set("Content-Type", "application/json");
            res.status(200).send(`Channel with id "${listingID}" was updated to "${JSON.stringify(result)}"`);
            db.close();
          }
        );
    }
  });
})
.delete((req, res) => {
  if (!req.get("X-User") || req.get("X-User").length == 0) {
    res.status(401).send("Unauthorized");
    return
  }
  if (req.params.id.length != 24) {
    res.status(403).send("Invalid listing ID");
    return
  }
  let listingID = new mongo.ObjectId(req.params.id);
  let reqUserID = JSON.parse(req.header("X-User")).id;
  MongoClient.connect(url, async (err, db) => {
    if (err) throw err;
    let dbo = db.db("mydb");
    const result = await dbo.collection("listings").findOne({ _id: listingID });
    if (!result) {
      res.status(403).send("Listing does not exist");
    } else {
      dbo
        .collection("listings")
        .findOneAndDelete(
          { _id: listingID },
          (err, document) => {
            if (err) throw err;
            let result = document.value;
            res.set("Content-Type", "application/json");
            res.json(result);
            db.close();
          }
        );
    }
  });
})

app.route("/v1/listings/creator/:id")
  .get((req, res) => {
    MongoClient.connect(url, (err, db) => {
      if (err) throw err;
      let dbo = db.db("mydb");
      dbo
        .collection("listings")
        .find({ creator: Number(req.params.id) })
        .toArray(function (err, result) {
          if (err) throw err;
          res.status(200);
          res.set(contentType, appJson);
          res.json(result);
          db.close();
        });
    })
  })

//start the server listening on host:port
app.listen(parseInt(port), () => {
  //callback is executed once server is listening
  console.log(`server is listening at http://${addr}...`);
});
