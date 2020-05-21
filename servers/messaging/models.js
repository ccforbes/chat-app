const mongoose = require("mongoose")

const UserSchema = mongoose.Schema({
    _id: false,
    id: Number,
    userName: String,
    firstName: String,
    lastName: String
})

const Channel = mongoose.model("Channel", new mongoose.Schema({
    name: { type: String, required: true, unique: true },
    description: String,
    private: Boolean,
    members: [UserSchema],
    creator: UserSchema
},
{ 
    timestamps: { createdAt: "createdAt", updatedAt: "editedAt"}
}))

const Message = mongoose.model("Message", new mongoose.Schema({
    channelID: { type: String, required: true },
    body: { type: String, required: true },
    creator: UserSchema
},
{
    timestamps: { createdAt: "createdAt", updatedAt: "editedAt"}
}))

const UserProfile = mongoose.model("User Profile", UserSchema)

module.exports = { UserProfile, Channel, Message }
