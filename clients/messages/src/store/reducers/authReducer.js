const initState = {
    user: null,
    expiration: null,
    error: null
}

const authReducer = (state = initState, action) => {
    switch (action.type) {
        case "login-success":
            return {
                ...state,
                user: {
                    userName: action.payload.user.userName,
                    firstName: action.payload.user.firstName,
                    lastName: action.payload.user.lastName
                },
                expiration: action.payload.expiration,
                error: null
            }
        case "login-error":
            return {
                ...state,
                error: action.payload
            }
        case "logout-success":
            return {
                ...state,
                user: null,
                expiration: null
            }
        case "update-success":
            return {
                ...state,
                user: {
                    userName: action.payload.userName,
                    firstName: action.payload.firstName,
                    lastName: action.payload.lastName
                },
                error: null
            }
        case "update-error":
            return {
                ...state,
                error: action.payload
            }
        case "auth-renewed":
            return {
                ...state,
                expiration: action.payload
            }
        case "auth-expired":
            return {
                user: null,
                expiration: null,
                error: null
            }
        default:
            return state
    }
}

export default authReducer