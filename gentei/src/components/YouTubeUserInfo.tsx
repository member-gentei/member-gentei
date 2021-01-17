import firebase from "firebase/app"
import "firebase/firestore"
import React, { useState, useEffect } from "react"
import { LoginStatus } from "../pages/Login"
import ChannelConnection from "./ChannelConnection"
import YouTubeMembershipCard, { YouTubeMembershipCardProps } from "./YouTubeMembershipCard"

interface iUserData {
    YoutubeChannelID?: string
    CandidateChannels: firebase.firestore.DocumentReference[]
    Memberships: firebase.firestore.DocumentReference[] | null
}

const YouTubeUserInfo = ({ user, loginStatus }: { user: firebase.User, loginStatus?: LoginStatus }) => {
    const [userData, setUserData] = useState<iUserData>()
    useEffect(() => {
        const fs = firebase.firestore()
        if (process.env.NODE_ENV === "development") {
            setUserData({
                YoutubeChannelID: "UCvpredjG93ifbCP1Y77JyFA",
                CandidateChannels: [
                    fs.doc("/channels/nakiri-ayame")
                ],
                Memberships: [
                    fs.doc("/channels/nakiri-ayame")
                ],
            })
        } else {
            return fs.collection("users").doc(user.uid).onSnapshot(doc => {
                setUserData(doc.data() as iUserData)
            })
        }
    }, [user])
    if (typeof userData === "undefined") {
        return null
    }
    // get memberships and candidate channels
    let memberMap: { [e: string]: boolean } = {}
    let cards: JSX.Element[] = []
    if (userData.Memberships != null) {
        userData.Memberships.forEach(e => {
            memberMap[e.path] = true
        })
        let memberOf: YouTubeMembershipCardProps[] = [], nonMember: YouTubeMembershipCardProps[] = []
        userData.CandidateChannels.forEach(e => {
            if (memberMap[e.path]) {
                memberOf.push({
                    path: e.path,
                    docRef: e,
                    isMember: true,
                })
            } else {
                nonMember.push({
                    path: e.path,
                    docRef: e,
                    isMember: false,
                })
            }
        })
        const pushCard = (e: YouTubeMembershipCardProps) => {
            cards.push(
                <YouTubeMembershipCard
                    key={e.path}
                    path={e.path}
                    docRef={e.docRef}
                    isMember={e.isMember} />
            )
        }
        memberOf.concat(nonMember)
        memberOf.forEach(pushCard)
        nonMember.forEach(pushCard)
    }
    return (
        <div className="connection">
            <ChannelConnection user={user} channelID={userData.YoutubeChannelID} loginStatus={loginStatus} />
            <div className="columns is-multiline mt-4">
                {cards}
            </div>
        </div>
    )
}

export default YouTubeUserInfo