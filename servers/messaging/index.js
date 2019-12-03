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
    // console.log(req.get("X-User"))
    // if (!req.get("X-User") || req.get("X-User").length == 0) {
    //   res.status(401).send("Unauthorized");
    //   return;
    // }
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
  console.log("req", req.headers)
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
      console.log("listingObj", listingObj)
      // console.log("result", result)
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
          console.log("result", result)
          console.log("req param id", req.params.id)
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
  console.log("req", req.headers)
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
    console.log(reqUserID, result.creator)
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


// app
//   .route("/v1/channels")
//   .get((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }

//     let reqUserID = JSON.parse(JSON.stringify(req.header("X-User"))).id;
//     MongoClient.connect(url, function (err, db) {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       dbo
//         .collection("channels")
//         .find(req.query.startsWith && { name: new RegExp("^" + req.query.startsWith) }, { creator: reqUserID })
//         .toArray(function (err, result) {
//           if (err) throw err;
//           res.status(200);
//           res.set(contentType, appJson);
//           res.json(result);
//           db.close();
//         });
//     });
//   })
//   .post((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     let reqUserID = JSON.parse(req.header("X-User")).id;
//     res.set("Content-Type", "application/json");
//     var myobj = req.body;
//     if (!myobj.name) {
//       res.status(403).send("Channel name is required");
//     } else {
//       MongoClient.connect(url, function (err, db) {
//         if (err) throw err;
//         let dbo = db.db("mydb");
//         if (myobj.name == null) myobj.name = "";
//         if (myobj.description == null) myobj.description = "";
//         if (myobj.private == null) myobj.private = false;
//         myobj.members == null
//           ? (myobj.members = [{ id: reqUserID }])
//           : (myobj.members = [{ id: reqUserID }].concat(myobj.members));
//         myobj.createdAt = new Date();
//         myobj.creator = reqUserID;
//         myobj.editedAt = null;
//         dbo.collection("channels").insertOne(myobj, function (err, result) {
//           if (err) throw err;
//           let channelMembers = [];
//           myobj.members && myobj.private && myobj.members.forEach(member => {
//             channelMembers.push(member.id);
//           });
//           const event = {
//             type: "channel-new",
//             channel: myobj,
//             userIDs: channelMembers
//           };
//           channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//           res.status(201);
//           res.set(contentType, appJson);
//           res.json(myobj);
//           db.close();
//         });
//       });
//     }
//   });

// app
//   .route("/v1/channels/:channelID")
//   .get((req, res) => {
//     res.set(contentType, appJson);
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let newChannelID = new mongo.ObjectId(req.params.channelID);
//     let reqUserID = JSON.parse(req.header("X-User")).id;
//     MongoClient.connect(url, function (err, db) {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       dbo.collection("channels").findOne(
//         {
//           _id: newChannelID
//         },
//         function (err, result) {
//           if (err) throw err;
//           let memberArr = function () {
//             for (let x of result.members) {
//               if (x.id == reqUserID) {
//                 return true;
//               }
//             }
//             return false;
//           };
//           if (result && result.private == true && !memberArr()) {
//             res.status(403).send("Current user is not member of channel " + newChannelID);
//           } else {
//             dbo
//               .collection("messages")
//               .find(req.query.before && { createdAt: { $lte: new Date(Number(req.query.before)) } }, {
//                 channelID: newChannelID
//               })
//               .sort({ createdAt: -1 })
//               .limit(100)
//               .toArray(function (err, result2) {
//                 if (err) throw err;

//                 res.send(result2);
//                 db.close();
//               });
//             db.close();
//           }
//         }
//       );
//     });
//   })
//   .post((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     let reqUserID = JSON.parse(req.header("X-User")).id;
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let channelID = new mongo.ObjectId(req.params.channelID);
//     MongoClient.connect(url, function (err, db) {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       dbo.collection("channels").findOne(
//         {
//           _id: channelID
//         },
//         function (err, result) {
//           if (err) throw err;
//           let memberArr = function () {
//             for (let x of result.members) {
//               if (x.id == reqUserID) {
//                 return true;
//               }
//             }
//             return false;
//           };
//           if (result == null || (result.private == true && !memberArr())) {
//             res.status(403).send("Current user is not member of channel " + channelID);
//           } else {
//             let msgObj = req.body;
//             msgObj.body = req.body.body;
//             msgObj.channelID = channelID;
//             msgObj.createdAt = new Date();
//             msgObj.creator = JSON.parse(req.get("X-User"));
//             msgObj.editedAt = null;
//             dbo.collection("messages").insertOne(msgObj, async function (err, result) {
//               if (err) throw err;
//               db.close();
//               let channelMembers = [];
//               const channel = await dbo.collection("channel").findOne({ _id: msgObj.channelID })
//               channel.members && channel.private && channel.members.forEach(member => {
//                 channelMembers.push(member.id);
//               });
//               const event = {
//                 type: "message-new",
//                 message: msgObj,
//                 userIDs: channelMembers
//               };
//               channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//               res.status(201);
//               res.set("Content-Type", "application/json");
//               res.json(msgObj);
//             });
//           }
//         }
//       );
//     });
//   })
//   .patch((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let channelID = new mongo.ObjectId(req.params.channelID);
//     let reqUserID = JSON.parse(req.header("X-User")).id;
//     let reqBody = req.body;
//     MongoClient.connect(url, async (err, db) => {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("channels").findOne({ _id: channelID });
//       if (!result || reqUserID != result.creator) {
//         res.status(403).send("Unauthorized user or channel does not exist");
//       } else {
//         dbo
//           .collection("channels")
//           .findOneAndUpdate(
//             { _id: channelID },
//             { $set: { name: reqBody.name, description: reqBody.description, editedAt: new Date() } },
//             { returnOriginal: false },
//             (err, document) => {
//               if (err) throw err;
//               let result = document.value;
//               let channelMembers = [];
//               result.members && result.private && result.members.forEach(member => {
//                 channelMembers.push(member.id);
//               });
//               const event = {
//                 type: "channel-update",
//                 channel: result,
//                 userIDs: channelMembers
//               };
//               channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//               res.set("Content-Type", "application/json");
//               res.status(200).send(`Channel with id "${channelID}" was updated to "${JSON.stringify(result)}"`);
//               db.close();
//             }
//           );
//       }
//     });
//   })
//   .delete((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let channelID = new mongo.ObjectId(req.params.channelID);
//     let reqUserID = JSON.parse(req.header("X-User")).id;
//     MongoClient.connect(url, async (err, db) => {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("channels").findOne({ _id: channelID });
//       if (!result || reqUserID != result.creator) {
//         res.status(403).send("Unauthorized user or channel does not exist");
//       } else {
//         dbo.collection("messages").remove({ channelID: channelID }, async (err, document) => {
//           if (err) throw err;
//           await dbo.collection("channels").remove({ _id: channelID });
//           let result = document.value;
//           let channelMembers = [];
//           result.members && result.private && result.members.forEach(member => {
//             channelMembers.push(member.id);
//           });
//           const event = {
//             type: "channel-delete",
//             channelID: result.id,
//             userIDs: channelMembers
//           };
//           channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//           res.set("Content-Type", "application/json");
//           res.status(200).send(`Channel with id "${channelID}" was successfuly deleted!"`);
//           db.close();
//         });
//       }
//     });
//   });

// app
//   .route("/v1/channels/:channelID/members")
//   .post((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length === 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let channelID = new mongo.ObjectId(req.params.channelID);
//     const currentUser = JSON.parse(req.header("X-User"));
//     const userBody = req.body;
//     MongoClient.connect(url, async (err, db) => {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("channels").findOne({ _id: channelID });
//       if (!result || JSON.stringify(result.creator) !== JSON.stringify(currentUser)) {
//         res.status(403).send("Unauthorized user or channelID does not exist");
//       } else {
//         dbo
//           .collection("channels")
//           .findOneAndUpdate(
//             { _id: channelID },
//             { $push: { members: userBody } },
//             { returnOriginal: false },
//             (err, document) => {
//               if (err) throw err;
//               const result = document.value;
//               res.status(201).send(`User with id "${userBody.id}" was added as a member to channel "${channelID}"`);
//               db.close();
//             }
//           );
//       }
//     });
//   })
//   .delete((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length === 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.channelID.length != 24) {
//       res.status(403).send("Invalid channel ID");
//     }
//     let channelID = new mongo.ObjectId(req.params.channelID);
//     const currentUser = JSON.parse(req.header("X-User"));
//     const userBody = req.body;
//     MongoClient.connect(url, async (err, db) => {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("channels").findOne({ _id: channelID });
//       if (!result || JSON.stringify(result.creator) !== JSON.stringify(currentUser)) {
//         res.status(403).send("Unauthorized user or channelID does not exist");
//       } else {
//         dbo
//           .collection("channels")
//           .findOneAndUpdate(
//             { _id: channelID },
//             { $pull: { members: userBody } },
//             { returnOriginal: false },
//             (err, document) => {
//               if (err) throw err;
//               res.status(200).send(`User with id "${userBody.id}" was removed as a member to channel "${channelID}"`);
//               db.close();
//             }
//           );
//       }
//     });
//   });

// app
//   .route("/v1/messages/:messageID")
//   .patch((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.messageID.length != 24) {
//       res.status(403).send("Invalid message ID");
//     }
//     let messageID = new mongo.ObjectId(req.params.messageID);
//     const user = JSON.parse(req.header("X-User"));
//     MongoClient.connect(url, async function (err, db) {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("messages").findOne({ _id: messageID });
//       if (!result || user.id != result.creator.id) {
//         res.status(403).send("Unauthorized user or messageID does not exist");
//       } else {
//         dbo
//           .collection("messages")
//           .findOneAndUpdate(
//             { _id: messageID },
//             { $set: { body: req.body.body, editedAt: new Date() } },
//             { returnOriginal: false },
//             (err, document) => {
//               if (err) throw err;
//               let result = document.value;
//               let channelMembers = [];
//               const channel = await dbo.collection("channel").findOne({ _id: result.channelID })
//               channel.members && channel.private && channel.members.forEach(member => {
//                 channelMembers.push(member.id);
//               });
//               const event = {
//                 type: "message-update",
//                 message: result,
//                 userIDs: channelMembers
//               };
//               channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//               res.set("Content-Type", "application/json");
//               res.status(200).send(`Message with id "${messageID}" was updated to "${JSON.stringify(result)}"`);
//               db.close();
//             }
//           );
//       }
//     });
//   })
//   .delete((req, res) => {
//     if (!req.get("X-User") || req.get("X-User").length == 0) {
//       res.status(401).send("Unauthorized");
//     }
//     if (req.params.messageID.length != 24) {
//       res.status(403).send("Invalid message ID");
//     }
//     let messageID = new mongo.ObjectId(req.params.messageID);
//     const currentUser = JSON.parse(req.header("X-User"));
//     MongoClient.connect(url, async (err, db) => {
//       if (err) throw err;
//       let dbo = db.db("mydb");
//       const result = await dbo.collection("messages").findOne({ _id: messageID });
//       if (!result || JSON.stringify(result.creator) != JSON.stringify(currentUser)) {
//         res.status(403).send("Unauthorized user or messageID does not exist");
//       } else {
//         dbo.collection("messages").remove({ _id: messageID }, (err, deleted) => {
//           if (err) throw err;
//           let channelMembers = [];
//           const channel = await dbo.collection("channel").findOne({ _id: result.channelID })
//           channel.members && channel.private && channel.members.forEach(member => {
//             channelMembers.push(member.id);
//           });
//           const event = {
//             type: "message-delete",
//             messageID: result.id,
//             userIDs: channelMembers
//           };
//           channel.sendToQueue(process.env.RABBITMQADDR, new Buffer(JSON.stringify(event)), { persistent: true });
//           res.status(200).send(`Message with ID ${messageID} was successful deleted!`);
//         });
//       }
//     });
//   });

//start the server listening on host:port
app.listen(parseInt(port), () => {
  //callback is executed once server is listening
  console.log(`server is listening at http://${addr}...`);
});
