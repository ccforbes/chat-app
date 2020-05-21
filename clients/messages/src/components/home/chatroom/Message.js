import React from 'react'

function Message(props) {
    const { message, showAuthor } = props
    const today = new Date()
    const msgDate = new Date(message.createdAt)
    return (
        <div className="message">
            {showAuthor ? <div className="message-username">
                {message.creator.userName} • {
                    today - msgDate < 86400000 ? 
                        msgDate.toLocaleTimeString() : msgDate.toLocaleDateString()
                }
            </div> : null}
            <div className="message-text">{message.body}</div>
        </div>
    )
}

export default Message