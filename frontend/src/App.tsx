import React from "react";
import "./App.scss";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import PrivacyPolicy from "./pages/PrivacyPolicy";
import AppIndex from "./pages/app/AppIndex";
import { LoginDiscord } from "./pages/Login";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/privacy" element={<PrivacyPolicy />} />
        <Route path="/login/discord" element={<LoginDiscord />} />
        <Route path="/app" element={<AppIndex />}></Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
