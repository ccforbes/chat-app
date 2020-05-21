import React, { useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Link, Redirect } from 'react-router-dom'
import { loginAction } from '../../store/actions/authActions'

function Login() {
    const [credentials, setCredentials] = useState({
        email: "",
        password: ""
    })
    const { user, error } = useSelector(state => state.auth)
    const dispatch = useDispatch()
    const login = (credentials) => dispatch(loginAction(credentials))

    let handleChange = (event) => {
        const { name, value } = event.target
        setCredentials({
            ...credentials,
            [name]: value
        })
    }

    let handleSubmit = (event) => {
        event.preventDefault()
        login(credentials)
    }

    if (user) {
        return <Redirect to="/home" />
    }

    return (
        <div className="container">
            <form className="white" onSubmit={handleSubmit}>
                <h5 className="grey-text text-darken-3">Login</h5>

                <div className="input-field">
                    <label className="active" htmlFor="email">Email</label>
                    <input 
                        type="email" 
                        name="email"
                        id="email" 
                        onChange={handleChange}
                    />
                </div>

                <div className="input-field">
                    <label className="active" htmlFor="password">Password</label>
                    <input 
                        type="password" 
                        name="password"
                        id="password" 
                        onChange={handleChange}
                    />
                </div>
                
                <div className="input-field">
                    <button className="btn lighten-1">Login</button>
                    <div className="red-text center">
                        { error ? <p>{error}</p> : null }
                    </div>
                </div>
                <Link to="/">Create an account</Link>
            </form>
        </div>
    )
}

export default Login