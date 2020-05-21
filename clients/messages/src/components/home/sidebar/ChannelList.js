import React from 'react'
import { getMessagesAction } from '../../../store/actions/channelActions'
import { useDispatch, useSelector } from 'react-redux'

function ChannelList(props) {
    const { channels } = props
    const dispatch = useDispatch()
    const { channelID } = useSelector(state => state.channel)
    const getMessages = (channel, roomID) => dispatch(getMessagesAction(channel, roomID))

    return (
        <div>
            <ul>
                <h6>Channels</h6>
                    {channels ? channels.map(channel => {
                        const active = channelID === channel._id ? "active" : ""
                        return (
                            <li key={channel._id} className={"room " + active}>
                                <a onClick={() => getMessages(channel, channel._id)} href={"#"+channel.name}>
                                    #{channel.name}
                                </a>
                            </li>
                        )
                    }) : null}
            </ul>
        </div>
    )
}

export default ChannelList