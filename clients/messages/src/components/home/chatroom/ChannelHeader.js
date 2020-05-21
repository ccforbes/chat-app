import React from 'react'

function ChannelHeader(props) {
    const { currChannel } = props
    return (
        <div className="channel-header">
            <h5>#{currChannel ? currChannel.name : null}</h5>
            <p>{currChannel ? currChannel.description : null}</p>
        </div>
    )
}

export default ChannelHeader