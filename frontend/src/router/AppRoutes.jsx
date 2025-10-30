import {BrowserRouter, Routes, Route} from "react-router-dom";

import DashboardLayout from "../layouts/DashboardLayout";

import Home from "../pages/home/Home";
import BotsList from "../pages/bots/BotsList";
import BotDetails from "../pages/bots/BotDetails";
import CreateBot from "../pages/createBot/CreateBot";
import LandingPage from "../pages/Landing/LandingPage";
import Signin from "../pages/auth/Register";
import Login from "../pages/auth/Login";

export default function AppRouter()
{
    return (
        <BrowserRouter>
            <Routes>

                {/* Landing Page (no layout) */}
                <Route path="/" element={<LandingPage />} />
                <Route path="/login" element={<Login />} />

                <Route path="/sign-up" element={<Signin />} />
                {/* All Dashboard routes inside layout */}
                <Route element={<DashboardLayout />}>
                    <Route path="/dashboard" element={<Home />} />
                    <Route path="/bots" element={<BotsList />} />
                    <Route path="/bots/:id" element={<BotDetails />} />
                    <Route path="/create" element={<CreateBot />} />
                </Route>

            </Routes>
        </BrowserRouter>
    );
}
