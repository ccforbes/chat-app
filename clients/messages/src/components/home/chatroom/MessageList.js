import React, { Component } from 'react'
import ReactDOM from 'react-dom'
import Message from './Message'

class MessageList extends Component {

    componentWillUpdate() {
        const node = ReactDOM.findDOMNode(this)
        this.shouldScrollToBottom = node.scrollTop + node.clientHeight + 100 >= node.scrollHeight
    }
    
    componentDidUpdate() {
        if (this.shouldScrollToBottom) {
            const node = ReactDOM.findDOMNode(this)
            node.scrollTop = node.scrollHeight   
        }
    }

    render() {
        let authorOfLastMessage = ""
        let timeOfLastMessage = new Date(null)
        let showAuthor = false
        return (
            <div className="message-list">
                {this.props.messages ? this.props.messages.map(message => {
                    const timeOfCurrMessage = new Date(message.createdAt)
                    const significantTimeDiff = timeOfCurrMessage - timeOfLastMessage > 120000
                    if (authorOfLastMessage !== message.creator.userName || significantTimeDiff) {
                        authorOfLastMessage = message.creator.userName
                        showAuthor = true
                    } else {
                        showAuthor = false
                    }
                    timeOfLastMessage = timeOfCurrMessage
                    return (
                        <Message key={message._id} message={message} showAuthor={showAuthor} />
                    )
                }) : null}
            </div>
        )
    }
}

export default MessageList