import firebase from 'firebase/app'
import "firebase/auth"
import "firebase/firestore"
import React from 'react'
import { LoginStatus } from '../pages/Login'
import LoginErrorDisplay from './LoginErrorDisplay'
import YouTubeUserInfo from './YouTubeUserInfo'

interface UserInfoProps {
    user?: firebase.User | null
    discordLoginStatus?: LoginStatus
    youTubeLoginStatus?: LoginStatus
}

const loginTest = async () => {
    try {
        await firebase.auth().createUserWithEmailAndPassword("mg-test@tindabox.net", "mg-test@tindabox.net")
    } catch (e) {
        if (e.code === "auth/email-already-in-use") {
            await firebase.auth().signInWithEmailAndPassword("mg-test@tindabox.net", "mg-test@tindabox.net")
        }
    }
    firebase.auth().currentUser?.updateProfile({
        displayName: "Test User",
    })
    firebase.auth().setPersistence(firebase.auth.Auth.Persistence.SESSION)
}

const UserInfo = ({ user, discordLoginStatus, youTubeLoginStatus }: UserInfoProps) => {
    const loginURL = encodeURIComponent(window.location.protocol + "//" + window.location.host + "/login/discord")
    const discordLoginURL = "https://discord.com/api/oauth2/authorize?client_id=768486576388177950&redirect_uri=" + loginURL + "&response_type=code&scope=identify%20guilds"
    let devLoginElement = null
    if (process.env.NODE_ENV === "development") {
        devLoginElement = (
            <button onClick={loginTest} className="button is-link">dev: trigger login</button>
        )
    }
    if (typeof user === 'undefined') {
        return (
            <div className="has-text-centered">
                <div className="lds-ring"><div></div><div></div><div></div><div></div></div>
                <p>Loading...</p>
            </div>
        )
    } else if (user === null) {
        return (
            <div className="has-text-centered">
                <div className="content">
                    <a href={discordLoginURL}><button className="button is-link">Sign in with Discord</button></a>
                    <LoginErrorDisplay loginStatus={discordLoginStatus} />
                    {devLoginElement}
                </div>
                <div className="content">
                    <p>
                        Sign in with Discord to get started! <br />
                        <em>Gentei</em> does not use the YouTube account you have connected in Discord, so you'll be associating one here separately.
                    </p>
                </div>
            </div >
        )
    }
    const signOut = () => {
        firebase.auth().signOut()
    }
    return (
        <div>
            <h2 className="title">Membership status</h2>
            <div className="content discord-info">
                <p>Hi, <strong>{user.displayName}</strong>.</p>
            </div>
            <YouTubeUserInfo user={user} loginStatus={youTubeLoginStatus} />
            <div className="buttons is-centered mt-6">
                <button className="button disabled">Request re-evaluation</button>
                <button onClick={signOut} className="button is-danger">Log out</button>
            </div>
        </div>
    )
}

export default UserInfo