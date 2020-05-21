import React from 'react'
import './App.css'
import { BrowserRouter, Switch, Route } from 'react-router-dom'
import Login from './components/auth/Login'
import Signup from './components/auth/Signup'
import Home from './components/home/Home'
import Update from './components/home/sidebar/update/Update'
import { useSelector, useDispatch } from 'react-redux'
import { authRenewed, authExpired } from './store/actions/authActions'

function App() {
    const { expiration } = useSelector(state => state.auth)
    const dispatch = useDispatch()
    const now = new Date()

    if (expiration && now.getTime() > expiration) {
        console.log("expired, logging out")
        dispatch(authExpired())
    } 

    return (
        <BrowserRouter>
            <div>
                <Switch>
                    <Route exact path="/" component={Signup} />
                    <Route path="/login" component={Login} />
                    <Route path="/home" component={Home} />
                    <Route path="/update" component={Update} />
                </Switch>
            </div>
        </BrowserRouter>
    )
}

export default App
