import React, { useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Link, Redirect } from 'react-router-dom'
import { signupAction } from '../../store/actions/authActions'

function Signup() {
    const [newUser, setNewUser] = useState({
        email: "",
        password: "",
        passwordConf: "",
        userName: "",
        firstName: "",
        lastName: ""
    })
    const { user, error } = useSelector(state => state.auth)
    const dispatch = useDispatch()
    const signup = (newUser) => dispatch(signupAction(newUser))

    let handleChange = (event) => {
        const { name, value } = event.target
        setNewUser({
            ...newUser,
            [name]: value
        })
    }

    let handleSubmit = (event) => {
        event.preventDefault()
        signup(newUser)
    }

    if (user) {
        return <Redirect to="/home" />
    }

    return (
        <div className="container">
            <form className="white col s12" onSubmit={handleSubmit}>
                <h5 className="grey-text text-darken-3">Sign Up</h5>

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

                <div className="row">
                    <div className="input-field col s6">
                        <label className="active" htmlFor="email">Email</label>
                        <input 
                            type="email"
                            name="email" 
                            id="email"
                            onChange={handleChange}
                        />
                    </div>
                    <div className="input-field col s6">
                        <label className="active" htmlFor="userName">User Name</label>
                        <input 
                            name="userName" 
                            id="userName"
                            onChange={handleChange} 
                        />
                    </div>
                </div>

                <div className="row">
                    <div className="input-field col s6">
                        <label className="active" htmlFor="password">Password</label>
                        <input 
                            type="password" 
                            name="password"
                            id="password"
                            onChange={handleChange} 
                        />
                        <span className="helper-text" data-error="wrong" cata-success="right">
                            MINIMUM: 6 characters
                        </span>
                    </div>
                    <div className="input-field col s6">
                        <label className="active" htmlFor="password">Confirm Your Password</label>
                        <input 
                            type="password" 
                            name="passwordConf"
                            id="passwordConf"
                            onChange={handleChange} 
                        />
                    </div>
                </div>

                <div className="input-field">
                    <button className="btn lighten-1">Sign Up</button>
                    <div className="red-text center">
                        { error ? <p>{error}</p> : null }
                    </div>
                </div>
                <Link to="/login">I'm already a member</Link>
            </form>
        </div>
    )
}

export default Signup