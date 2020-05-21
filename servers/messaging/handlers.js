const ChannelsHandler = (req, res, { Channel, UserProfile, getRabbitChannel }) => {
    const user = JSON.parse(req.get("X-User"))
    switch (req.method) {
        case "GET":
            Channel.find({
                $or: [
                    { members: { $elemMatch: { id: user.id } } },
                    { private: { $ne: true } }
                ]
            })
                .exec()
                .then(docs => {
                    res.status(200).json(docs)
                })
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            return
        case "POST":
            const { name, description, private, members } = req.body
            const { id, userName, firstName, lastName } = user
            if (!name) {
                res.status(400).send("The channel name can not be empty.")
                return
            }
            const authUser = {
                id,
                userName,
                firstName,
                lastName
            }
            const creator = new UserProfile(authUser)
            let allMembers = []
            let userIDs = []
            if (private === true) {
                allMembers.push(creator)
                for (const member of members) {
                    const newMember = new UserProfile(member)
                    allMembers.push(newMember)
                    userIDs.push(member.id)
                }
            }
            const channel = new Channel({
                name: name,
                description: description,
                private: private,
                members: allMembers,
                creator: creator
            })
            channel.save()
                .then(result => {
                    let ch = getRabbitChannel()
                    ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                        {
                            type: "channel-new",
                            channel: result,
                            userIDs: userIDs
                        }
                    )))
                    res.status(201).json({
                        result
                    })
                })
                .catch(err => {
                    res.status(500).send("The channel could not be added.")
                })
            return  
        default:
            res.status(405).send("The request method is not allowed.")
            return
    }
}

const SpecificChannelsHandler = async (req, res, { Channel, Message, UserProfile, getRabbitChannel }) => {
    let ch = getRabbitChannel()
    const user = JSON.parse(req.get("X-User"))
    const chanID = req.params.channelID
    const channel = await Channel.find({
        $and: [
            { _id: chanID },
            { $or: [
                { private: { $ne: true } },
                { members: { $elemMatch: { id: user.id } } }
            ] }
        ]
    })
    const isPrivate = channel.private == true
    const userIDs = isPrivate ? channel.members.map(({id}) => id) : []

    if (channel.length == 0) {
        res.status(403).send("You're not a member of this channel")
        return
    }
    let isCreator = channel[0].creator.id == user.id

    switch (req.method) {
        case "GET":
            Message.find({channelID: chanID})
                .sort({ createdAt: -1 })
                .limit(100)
                .exec()
                .then(result => {
                    res.json(result)
                })
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            return
        case "POST":
            // create new message in channel using request body's JSON (just read body)
            const { body } = req.body
            const { id, userName, firstName, lastName } = user
            const creator = new UserProfile({
                id,
                userName,
                firstName,
                lastName
            })
            const message = new Message({
                channelID: chanID,
                body: body,
                creator: creator
            })
            message.save()
                .then(result => {
                    ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                        {
                            type: "message-new",
                            message: result,
                            userIDs: userIDs
                        }
                    )))
                    res.status(201).json(result)
                })
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            return
        case "PATCH":
            if (!isCreator) {
                res.status(403).send("You are not the creator of the channel")
                return
            }
            try {
                const updatedChannel = await Channel.findOneAndUpdate(
                    { _id: chanID }, 
                    { $set: req.body }, 
                    { new: true }
                )
                ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                    {
                        type: "channel-update",
                        channel: updatedChannel,
                        userIDs: userIDs
                    }
                )))
                res.json(updatedChannel)
            } catch(e) {
                res.status(500).json({
                    error: e
                })
            }
            return
        case "DELETE":
            if (!isCreator) {
                res.status(403).send("You are not the creator of the channel")
                return
            }
            Message.deleteMany({ channelID: chanID })
                .exec()
                .then()
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            Channel.deleteOne({ _id: chanID })
                .exec()
                .then(() => {
                    ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                        {
                            type: "channel-delete",
                            channelID: chanID,
                            userIDs: userIDs
                        }
                    )))
                    res.send("The channel has been deleted")
                })
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            
            return
        default:
            res.status(405).send("The request method is not allowed.")
            return
    }
}

const MembersHandler = async (req, res, { Channel, UserProfile }) => {
    const user = JSON.parse(req.get("X-User"))
    const channel = await Channel.findById(req.params.channelID)
    const isCreator = channel.creator.id == user.id
    if (!isCreator) {
        res.status(403).send("You're not the creator of the channel")
        return
    }
    if (req.body.id == "") {
        res.status(400).send("The user's ID is required")
        return
    }
    const chanID = req.params.channelID
    switch (req.method) {
        case "POST":
            const newMember = new UserProfile(req.body)
            channel.update({ $push: { members: newMember } })
                .exec()
                .then(res.status(201).send("Member added to channel " + chanID))
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            return
        case "DELETE":
            const member = new UserProfile(req.body)
            Channel.update ({ _id: chanID }, { $pull: { members: member } })
                .exec()
                .then(res.status(200).send("Member was removed from channel " + chanID))
            return
        default:
            res.status(405).send("The request method is not allowed.")
            return
    }
}

const MessagesHandler = async (req, res, { Channel, Message, getRabbitChannel }) => {
    let ch = getRabbitChannel()
    const user = JSON.parse(req.get("X-User"))
    const msgID = req.params.messageID
    const message = await Message.findById(req.params.messageID)
    const isCreator = message.creator.id == user.id
    const channel = await Channel.findById(message.channelID)
    const isPrivate = channel.private == true
    const userIDs = isPrivate ? channel.members.map(({id}) => id) : []
    if (!isCreator) {
        res.status(403).send("You're not the creator of the message")
        return
    }
    switch (req.method) {
        case "PATCH":
            const { body } = req.body
            try {
                let updatedMessage = await Message.findOneAndUpdate(
                    { _id: msgID }, 
                    { $set: { body: body } },
                    { new: true }
                )
                ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                    {
                        type: "message-update",
                        message: updatedMessage,
                        userIDs: userIDs
                    }
                )))
                res.json(updatedMessage)
            } catch(e) {
                res.status(500).json({
                    error: e
                })
            }
            return
        case "DELETE":
            Message.deleteOne({ _id: msgID })
                .exec()
                .then(() => {
                    ch.sendToQueue("msgs_queue", Buffer.from(JSON.stringify(
                        {
                            type: "message-delete",
                            messageID: msgID,
                            userIDs: userIDs
                        }
                    )))
                    res.send("Message has been deleted")
                })
                .catch(err => {
                    res.status(500).json({
                        error: err
                    })
                })
            return
        default:
            res.status(405).send("The request method is not allowed.")
            return
    }
}

module.exports = { 
    ChannelsHandler, 
    SpecificChannelsHandler, 
    MembersHandler, 
    MessagesHandler 
}