export const getChannelsAction = () => {
    return (dispatch, getState) => {
        fetch("https://api.bopboyz222.xyz/v1/channels", {
            headers: {
                "Content-Type": "application/json",
                "Authorization":localStorage.getItem("authToken")
            }
        }).then(resp => {
            if (!resp.ok) {
                throw Error(resp.statusText)
            }
            return resp.json
        }).then(json => {
            dispatch({
                type: "get-channels",
                payload: json
            })
        }).catch(err => {
            dispatch({
                type: "channels-error",
                payload: "Error receiving channels."
            })
        })
    }
}

export const getMessagesAction = (channel, roomID) => {
    return (dispatch, getState) => {
        fetch("https://api.bopboyz222.xyz/v1/channels/" + roomID, {
            headers: {
                "Content-Type": "application/json",
                "Authorization":localStorage.getItem("authToken")
            }
        }).then(resp => {
            if (!resp.ok) {
                throw Error(resp.statusText)
            }
            return resp.json()
        }).then(json => {
            const messages = json.reverse()
            dispatch({
                type: "get-messages",
                payload: {
                    messages: messages,
                    roomID: roomID,
                    currChannel: channel
                }
            })
        }).catch(err => {
            dispatch({
                type: "channels-error",
                payload: "Error receiving channels."
            })
        })
    }
}

export const sendMessageAction = (message) => {
    return (dispatch, getState) => {
        const channelID = getState().channel.channelID
        fetch("https://api.bopboyz222.xyz/v1/channels/" + channelID, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization":localStorage.getItem("authToken")
            },
            body: JSON.stringify({
                body: message
            })
        }).then(resp => {
            if (!resp.ok) {
                resp.text().then(text => {
                    throw Error(text)
                })
            } 
        }).catch(err => {
            dispatch({
                type: "message-error",
                payload: err
            })
        })
    }
}