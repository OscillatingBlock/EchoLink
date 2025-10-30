import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { connectTwilio, getMyNumber } from "./authAPI";

export const registerAndConnect = createAsyncThunk(
  "auth/registerAndConnect",
  async (payload, { rejectWithValue }) => {
    try {
      const data = await connectTwilio(payload);
      if (data?.access_token) {
        localStorage.setItem("access_token", data.access_token);
      }
      return data;
    } catch (e) {
      return rejectWithValue(e.message || "Registration failed");
    }
  }
);

export const fetchMyNumber = createAsyncThunk(
  "auth/fetchMyNumber",
  async (_, { rejectWithValue }) => {
    try {
      const data = await getMyNumber();
      return data;
    } catch (e) {
      return rejectWithValue(e.message || "Failed to fetch number");
    }
  }
);

const authSlice = createSlice({
  name: "auth",
  initialState: {
    access_token: localStorage.getItem("access_token") || null,
    phone_number: null,
    bots_count: null,
    loading: false,
    error: null,
    message: null,
  },
  reducers: {
    logout(state) {
      state.access_token = null;
      state.phone_number = null;
      state.bots_count = null;
      localStorage.removeItem("access_token");
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(registerAndConnect.pending, (state) => {
        state.loading = true;
        state.error = null;
        state.message = null;
      })
      .addCase(registerAndConnect.fulfilled, (state, action) => {
        state.loading = false;
        state.message = action.payload?.message || "Connected";
        state.access_token = action.payload?.access_token || null;
        state.phone_number = action.payload?.phone_number || null;
      })
      .addCase(registerAndConnect.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Registration error";
      })
      .addCase(fetchMyNumber.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMyNumber.fulfilled, (state, action) => {
        state.loading = false;
        state.phone_number = action.payload?.phone_number || state.phone_number;
        state.bots_count = action.payload?.bots_count ?? state.bots_count;
      })
      .addCase(fetchMyNumber.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || "Failed to fetch number";
      });
  },
});

export const { logout } = authSlice.actions;
export default authSlice.reducer;
