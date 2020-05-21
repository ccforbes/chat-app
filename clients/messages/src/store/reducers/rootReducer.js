import authReducer from './authReducer'
import channelReducer from './channelReducer'
import { combineReducers } from 'redux'

const rootReducer = combineReducers({
    auth: authReducer,
    channel: channelReducer
})

export default rootReducer