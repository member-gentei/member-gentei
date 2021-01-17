import firebase from 'firebase/app';
import "firebase/auth"
import { useState, useEffect } from "react"
import Footer from "../components/Footer"
import UserInfo from '../components/UserInfo';
import { LoginStatus } from './Login';

interface MembershipManagerProps {
    ready: boolean
    discordLoginStatus?: LoginStatus
    youTubeLoginStatus?: LoginStatus
}

const MembershipManager = ({ ready, discordLoginStatus, youTubeLoginStatus }: MembershipManagerProps) => {
    const [user, setUser] = useState<firebase.User | null | undefined>(undefined)
    useEffect(() => {
        if (ready) {
            return firebase.auth().onAuthStateChanged(fsUser => {
                setUser(fsUser)
            })
        }
    }, [ready])
    return (
        <div>
            <section className="hero is-primary">
                <div className="hero-body">
                    <div className="container">
                        <h1 className="title">Gentei / 限定</h1>
                        <h2 className="subtitle">VTuber channel membership verification</h2>
                    </div>
                </div>
            </section>
            <section role="main" className="section">
                <div className="container">
                    <UserInfo user={user} discordLoginStatus={discordLoginStatus} youTubeLoginStatus={youTubeLoginStatus} />
                </div>
            </section>
            <section id="qa" className="section">
                <div className="container">
                    <h2 className="title is-4">Q&amp;A</h2>
                    <h3 className="subtitle is-5 mt-4">How does this work?</h3>
                    <p><em>Gentei</em> is connected to designated bots in VTuber fan Discord servers. It uses the YouTube API to fetch channel information in an more roundabout - but similarly effective - manner as Discord's official <a href="https://support.discord.com/hc/en-us/articles/215162978-Youtube-Channel-Memberships-Integration-FAQ">YouTube channel memberships integration</a>.</p>
                    <br />
                    <p><em>Gentei</em> periodically checks your membership and notifies the Discord bots of membership status changes.</p>
                </div>
            </section>
            <Footer withYouTubeImage={true} />
        </div>
    )
}

export default MembershipManager