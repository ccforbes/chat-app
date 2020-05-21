export const channelAction = (socketMsgData) => {
    return (dispatch, getState) => {
        dispatch({
            type: socketMsgData.type,
            payload: socketMsgData.channel
        })
    }
}

export const messageAction = (socketMsgData) => {
    return (dispatch, getState) => {
        dispatch({
            type: socketMsgData.type,
            payload: socketMsgData.message
        })
    }
}