import {configureStore} from "@reduxjs/toolkit";
import botsReducer from "../features/bots/botsSlice";
import authReducer from "../features/auth/authSlice";

export const store=configureStore({
    reducer: {
        bots: botsReducer,
        auth: authReducer,
    },
});

export default store;
