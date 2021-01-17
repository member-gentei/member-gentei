import firebase from 'firebase/app';
import "firebase/auth"
import "firebase/firestore"
import React, { useCallback, useState } from 'react';
import { useEffect } from "react"
import { BrowserRouter, Route, Switch, useHistory, useLocation } from "react-router-dom";
import './App.css';
import Home from './pages/Home';
import Login, { LoginStatus } from './pages/Login';
import MembershipManager from './pages/MembershipManager';

const firebaseConfig = {
  apiKey: "AIzaSyAf_lA-srGXkaSrM2mb-py9c4VMa0zJcSY",
  authDomain: "member-gentei.firebaseapp.com",
  databaseURL: "https://member-gentei.firebaseio.com",
  projectId: "member-gentei",
  storageBucket: "member-gentei.appspot.com",
  messagingSenderId: "649732146530",
  appId: "1:649732146530:web:68911af59aff1fd012183b"
}

function App() {
  const [firebaseReady, setFirebaseReady] = useState(false)
  const [loginStatus, setLoginStatus] = useState<LoginStatus>()
  useEffect(() => {
    (async () => {
      if (!firebaseReady) {
        firebase.initializeApp(firebaseConfig)
        if (process.env.NODE_ENV === "development") {
          console.debug("using Firebase emulators")
          firebase.auth().useEmulator("http://localhost:9099")
          firebase.firestore().useEmulator("localhost", 8099)
        }
      }
      setFirebaseReady(true)
    })()
  }, [firebaseReady])

  const onDiscordLogin = useCallback((status: LoginStatus) => {
    setLoginStatus(status)
  }, [setLoginStatus])
  return (
    <BrowserRouter>
      <Switch>
        <Route path="/app">
          <MembershipManager ready={firebaseReady} loginStatus={loginStatus} />
        </Route>
        <Route path="/login/:loginType">
          <Login
            onDiscordLogin={onDiscordLogin} />
        </Route>
        <Route path="/">
          <Home />
        </Route>
      </Switch>
    </BrowserRouter>
  )
}

export default App;
