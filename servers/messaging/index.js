"use strict"

const express = require("express")
const mongoose = require("mongoose")
const { UserProfile, Channel, Message } = require("./models")
const { 
    ChannelsHandler,
    SpecificChannelsHandler,
    MembersHandler,
    MessagesHandler
} = require("./handlers")
const amqp = require("amqplib/callback_api")

const port = process.env.PORT || 5001
const mongoEndpoint = "mongodb://mongocontainer:27017/messaging"

const rabbAddr = process.env.RABBITADDR || "amqp://guest:guest@rabbit:5672"
let rabbitChannel;

const getRabbitChannel = () => {
    return rabbitChannel
}

const app = express()
app.use(express.json())

const connect = () => {
    mongoose.connect(mongoEndpoint, {reconnectTries: 5})
}

const RequestWrapper = (handler, SchemeAndDbForwarder) => {
    return (req, res) => {
        const value = req.get("X-User")
        if (!value) {
            res.status(401).send("The user is not authenticated")
            return
        }
        handler(req, res, SchemeAndDbForwarder)
    }
}

app.use("/v1/channels/:channelID/members", RequestWrapper(MembersHandler, { Channel, UserProfile, getRabbitChannel } ))
app.use("/v1/channels/:channelID", RequestWrapper(SpecificChannelsHandler, { Channel, Message, UserProfile, getRabbitChannel }))
app.use("/v1/channels", RequestWrapper(ChannelsHandler, { Channel, UserProfile, getRabbitChannel }))
app.use("/v1/messages/:messageID", RequestWrapper(MessagesHandler, { Channel, Message, getRabbitChannel }))

connect()
mongoose.connection.on("error", console.error)
    .on("disconnected", connect)
    .once("open", main)

async function main() {
    amqp.connect(rabbAddr, (err, conn) => {
        if (err) {
            console.log("Error connecting to rabbit instance")
            process.exit(1)
        }

        conn.createChannel((err, ch) => {
            if (err) {
                console.log("Error creating channel")
                process.exit(1)
            }

            ch.assertQueue("msgs_queue", {durable: true})
            rabbitChannel = ch
        })
        app.listen(port, "", () => {
            console.log(`server listening ${port}...`)
        })
    })

}