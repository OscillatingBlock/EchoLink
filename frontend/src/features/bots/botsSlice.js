import {createSlice, createAsyncThunk} from "@reduxjs/toolkit";
import {fetchBots, createBot, fetchBotDetails, deleteBot} from "./botsAPI";

// LOAD ALL BOTS
export const getBots=createAsyncThunk("bots/getBots", async () =>
{
    return await fetchBots();
});

// CREATE NEW BOT
export const addBot=createAsyncThunk("bots/addBot", async (data) =>
{
    return await createBot(data);
});

// LOAD ONE BOT
export const getBotDetails=createAsyncThunk("bots/getBotDetails", async (id) =>
{
    return await fetchBotDetails(id);
});

// DELETE BOT
export const removeBot=createAsyncThunk("bots/removeBot", async (id) =>
{
    await deleteBot(id);
    return id; // return id to remove from state
});

const botsSlice=createSlice({
    name: "bots",
    initialState: {
        list: [],
        selectedBot: null,
        loading: false,
        error: null,
    },
    reducers: {},

    extraReducers: (builder) =>
    {
        builder
            // GET ALL
            .addCase(getBots.pending, (state) =>
            {
                state.loading=true;
            })
            .addCase(getBots.fulfilled, (state, action) =>
            {
                state.loading=false;
                state.list=action.payload;
            })
            .addCase(getBots.rejected, (state) =>
            {
                state.loading=false;
                state.error="Failed to load bots";
            })

            // CREATE
            .addCase(addBot.pending, (state) =>
            {
                state.loading=true;
            })
            .addCase(addBot.fulfilled, (state, action) =>
            {
                state.loading=false;
                state.list.push(action.payload);
            })

            // DETAILS
            .addCase(getBotDetails.pending, (state) =>
            {
                state.loading=true;
            })
            .addCase(getBotDetails.fulfilled, (state, action) =>
            {
                state.loading=false;
                state.selectedBot=action.payload;
            })

            // DELETE
            .addCase(removeBot.fulfilled, (state, action) =>
            {
                state.list=state.list.filter(bot => bot.id!==action.payload);
            });
    },
});

export default botsSlice.reducer;
