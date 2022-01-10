import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.scss";
import AppIndex from "./pages/app/AppIndex";
import Enrollment from "./pages/app/Enrollment";
import GuildAdmin from "./pages/app/GuildAdmin";
import UserDashboard from "./pages/app/UserDashboard";
import Home from "./pages/Home";
import { LoginDiscord, LoginYouTube } from "./pages/Login";
import PrivacyPolicy from "./pages/PrivacyPolicy";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/privacy" element={<PrivacyPolicy />} />
        <Route path="/login/discord" element={<LoginDiscord />} />
        <Route path="/login/youtube" element={<LoginYouTube />} />
        <Route path="/app" element={<AppIndex />}>
          <Route index element={<UserDashboard />} />
          <Route path="enroll" element={<Enrollment />} />
          <Route path="server/:guildID" element={<GuildAdmin />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
