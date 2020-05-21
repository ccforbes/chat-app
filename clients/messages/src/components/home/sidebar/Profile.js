import React from 'react'
import { useDispatch } from 'react-redux'
import { Link } from 'react-router-dom'
import { logoutAction } from '../../../store/actions/authActions'

function Profile(props) {
    const dispatch = useDispatch()
    const logout = () => dispatch(logoutAction())

    return (
        <div className="center">
            <h4>Welcome, {props.user.firstName} {props.user.lastName}!</h4>
            <Link to="/update"><button className="btn-small">Update</button></Link>
            <button onClick={logout} className="btn-small red">Logout</button>
        </div>
    )
}

export default Profile