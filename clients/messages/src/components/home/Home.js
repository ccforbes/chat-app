import React, { useEffect } from 'react'
import { Redirect } from 'react-router-dom'
import Sidebar from './sidebar/Sidebar'
import ChannelHeader from './chatroom/ChannelHeader'
import MessageList from './chatroom/MessageList'
import SendMessageForm from './chatroom/SendMessageForm'
import { useSelector, useDispatch } from 'react-redux'
import { channelAction, messageAction } from '../../store/actions/socketActions'

function Home() {
    const { user } = useSelector(state => state.auth)
    const { channels, messages, currChannel } = useSelector(state => state.channel)
    const dispatch = useDispatch()
    
    useEffect(() => {
        const authToken = localStorage.getItem("authToken")
        const wsApiEndpoint = "wss://api.bopboyz222.xyz/v1/ws?auth=" + authToken
        let socket = new WebSocket(wsApiEndpoint)
    
        socket.onopen = () => {
            console.log("Connection Opened")
        }
    
        socket.onclose = () => {
            console.log("Connection Closed")
        }
    
        socket.onmessage = msg => {
            console.log("Message received ")
            const data = JSON.parse(msg.data)
            if (data.type.includes("message")) {
                dispatch(messageAction(data))
            }
            if (data.type.includes("channel")) {
                dispatch(channelAction(data))
            }
        }
        return () => {
            socket.close()
        }
    }, [])

    if (!user) {
        return <Redirect to="/login" />
    }


    return (
        <div className="app">
            <Sidebar user={user} channels={channels}/>
            <ChannelHeader currChannel={currChannel}/>
            <MessageList messages={messages} />
            <SendMessageForm />
        </div>
    )
}

export default Home