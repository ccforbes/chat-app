import React from 'react'
import Profile from './Profile'
import ChannelList from './ChannelList'

function Sidebar(props) {
    return (
        <div className="rooms-list">
            <Profile user={props.user} />
            <ChannelList channels={props.channels}/>
        </div>
    )
}

export default Sidebar