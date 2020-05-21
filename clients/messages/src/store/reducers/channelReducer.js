const initState = {
    channels: null,
    currChannel: null,
    channelID: "",
    messages: null,
    error: null
}

const channelReducer = (state = initState, action) => {
    switch (action.type) {
        case "get-channels":
            return {
                ...state,
                channels: action.payload,
                error: null
            }
        case "channels-error":
            return {
                ...state,
                error: action.payload
            }
        case "remove-channels":
            return {
                ...state,
                channels: null,
                currChannel: null,
                channelID: "",
                messages: null
            }
        case "get-messages":
            return {
                ...state,
                channelID: action.payload.roomID,
                messages: action.payload.messages,
                currChannel: action.payload.currChannel
            }
        case "message-new":
            const socketMessages = [...state.messages, action.payload]
            return {
                ...state,
                messages: socketMessages
            }
        case "channel-new": 
            const updateChannels = [...state.channels, action.payload]
            return {
                ...state,
                channels: updateChannels
            }
        default:
            return state
    }
}

export default channelReducer