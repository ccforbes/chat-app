import React, { useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Link, Redirect } from 'react-router-dom'
import { updateAction } from '../../../../store/actions/authActions'

function Update(props) {
    const [userInfo, setUserInfo] = useState({
        firstName: "",
        lastName: ""
    })
    const { user, error } = useSelector(state => state.auth)
    const dispatch = useDispatch()
    const update = async (userInfo) => dispatch(updateAction(userInfo))

    let handleChange = (event) => {
        const { name, value } = event.target
        setUserInfo({
            ...userInfo,
            [name]: value
        })

    }

    let handleSubmit = async (event) => {
        event.preventDefault()
        await update(userInfo)
        props.history.push("/home")
    }

    if (!user) {
        return <Redirect to="/login" />
    }

    return (
        <div className="container">
            <form className="white col s12" onSubmit={handleSubmit}>
                <h5 className="grey-text text-darken-3">Update User Information</h5>

                <div className="row">
                    <div className="input-field col s6">
                        <label className="active" htmlFor="firstName">First Name</label>
                        <input 
                            name="firstName" 
                            id="firstName"
                            onChange={handleChange} 
                        />
                    </div>
                    <div className="input-field col s6">
                        <label className="active" htmlFor="lastName">Last Name</label>
                        <input 
                            name="lastName" 
                            id="lastName"
                            onChange={handleChange}
                        />
                    </div>
                </div>

                <div className="input-field">
                    <button className="btn">Update</button>
                    <div className="red-text center">
                        { error ? <p>{error}</p> : null }
                    </div>
                </div>
                <Link to="/home">Actually, never mind.</Link>
            </form>
        </div>
    )
}

export default Update