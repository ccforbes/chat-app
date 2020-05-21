import React, { useState } from 'react'
import { useDispatch } from 'react-redux'
import { sendMessageAction } from '../../../store/actions/channelActions'

function SendMessageForm() {
    const [message, setMessage] = useState("")
    const dispatch = useDispatch()
    const sendMessage = (message) => dispatch(sendMessageAction(message))

    let handleChange = (event) => {
        const { value } = event.target
        setMessage(value)
    }

    let handleSubmit = (event) => {
        event.preventDefault()
        sendMessage(message)
        setMessage("")
    }

    return (
        <form 
            className="send-message-form"
            onSubmit={handleSubmit}>
            <input  
                placeholder="Type your message and hit ENTER" 
                value={message}
                onChange={handleChange} 
            />
        </form>
    )
}

export default SendMessageForm