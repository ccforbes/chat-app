export const loginAction = (credentials) => {
    const newExpiration = getNewExpirationTime()
    return async (dispatch, getState) => {
        await fetch("https://api.bopboyz222.xyz/v1/sessions", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: credentials.email,
                password: credentials.password,
            })
        }).then(async resp => {
            if (!resp.ok) {
                localStorage.setItem("authToken", "")
                await resp.text().then(text => {
                    throw Error(text)
                })
            } else {
                localStorage.setItem("authToken", resp.headers.get("Authorization"))
                return resp.json() 
            }
        }).then(json => {
            dispatch({ 
                type: "login-success",
                payload: {
                    user: json,
                    expiration: newExpiration
                }
            })
        }).catch(err => {
            dispatch({
                type: "login-error",
                payload: err.message
            })
        })

        initializeChannelsAndMessages(dispatch)
    }
}

export const logoutAction = () => {
    return (dispatch, getState) => {
        fetch("https://api.bopboyz222.xyz/v1/sessions/mine", {
            method: "DELETE",
            headers: { 
                "Authorization":localStorage.getItem("authToken") 
            }
        }).then(resp => {
            if (resp.ok) {
                dispatch({
                    type: "logout-success"
                })
                dispatch({
                    type: "remove-channels"
                })
                localStorage.removeItem("authToken")
            }
        })
    }
}

export const signupAction = (newUser) => {
    const newExpiration = getNewExpirationTime()
    return async (dispatch, getState) => {
        await fetch("https://api.bopboyz222.xyz/v1/users", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: newUser.email,
                password: newUser.password,
                passwordConf: newUser.passwordConf,
                userName: newUser.userName,
                firstName: newUser.firstName,
                lastName: newUser.lastName
            })
        }).then(async resp => {
            if (!resp.ok) {
                localStorage.setItem("authToken", "")
                await resp.text().then(text => {
                    throw Error(text)
                })
            }
            localStorage.setItem("authToken", resp.headers.get("Authorization"))
            return resp.json()
        }).then(json => {
            dispatch({
                type: "login-success",
                payload: {
                    user: json,
                    expiration: newExpiration
                }
            })
        }).catch(err => {
            dispatch({
                type: "login-error",
                payload: err.message
            })
        })

        initializeChannelsAndMessages(dispatch)
    }
}

export const updateAction = (userInfo) => {
    return (dispatch, getState) => {
        fetch("https://api.bopboyz222.xyz/v1/users/me", {
            method: "PATCH",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("authToken")
            },
            body: JSON.stringify({
                firstName: userInfo.firstName,
                lastName: userInfo.lastName,
            })
        }).then(async resp => {
            if (!resp.ok) {
                await resp.text().then(text => {
                    throw Error(text)
                })
            }
            return resp.json()
        }).then(json => {
            dispatch({
                type: "update-success",
                payload: json
            })
        }).catch(err => {
            dispatch({
                type: "login-error",
                payload: err.message
            })
        })
    }
}

export const authRenewed = () => {
    const newExpiration = getNewExpirationTime()
    return (dispatch, getState) => {
        dispatch({
            type: "auth-renewed",
            payload: newExpiration
        })
    }
}

export const authExpired = () => {
    return (dispatch, getState) => {
        dispatch({
            type: "auth-expired"
        })
    }
}

function getNewExpirationTime() {
    const newExpiration = new Date()
    return newExpiration.getTime() + (60*60*1000)
}

async function initializeChannelsAndMessages(dispatch) {
    let general = null
    let generalID = "";

    // grab user's channels
    await fetch("https://api.bopboyz222.xyz/v1/channels", {
        headers: {
            "Content-Type": "application/json",
            "Authorization":localStorage.getItem("authToken")
        }
    }).then(resp => {
        if (!resp.ok) {
            throw Error(resp.statusText)
        }
        return resp.json()
    }).then(json => {
        general = json[0]
        generalID = general._id
        dispatch({
            type: "get-channels",
            payload: json
        })
    }).catch(err => {
        dispatch({
            type: "channels-error",
            payload: "Error receiving channels."
        })
    })

    // grab general messages
    fetch("https://api.bopboyz222.xyz/v1/channels/" + generalID, {
        headers: {
            "Content-Type": "application/json",
            "Authorization":localStorage.getItem("authToken")
        }
    }).then(resp => {
        if (!resp.ok) {
            throw Error(resp.statusText)
        }
        return resp.json()
    }).then(json => {
        const messages = json.reverse()
        dispatch({
            type: "get-messages",
            payload: {
                messages: messages, 
                roomID: generalID,
                currChannel: general
            }
        })
    }).catch(err => {
        dispatch({
            type: "channels-error",
            payload: "Error receiving channels."
        })
    })
}