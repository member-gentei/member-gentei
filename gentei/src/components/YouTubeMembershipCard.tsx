import firebase from 'firebase/app'
import "firebase/firestore"
import { useEffect, useState } from "react"


export interface YouTubeMembershipCardProps {
    path: string
    docRef: firebase.firestore.DocumentReference
    isMember: boolean
}

interface channelCardDoc {
    ChannelID: string
    ChannelTitle: string
    Thumbnail: string
}

const YouTubeMembershipCard = ({ path, docRef, isMember }: YouTubeMembershipCardProps) => {
    const [doc, setDoc] = useState<channelCardDoc>()
    useEffect(() => {
        (async () => {
            const data = (await docRef.get()).data()
            setDoc(data as channelCardDoc)
        })()
    }, [docRef])
    if (typeof doc === "undefined") {
        return (
            <div className="column is-one-quarter">
                <div className="card">
                    <div className="card-content">
                        <div className="lds-ring"><div></div><div></div><div></div><div></div></div>
                    </div>
                </div>
            </div>
        )
    }
    const channelHref = `https://youtube.com/channel/${doc.ChannelID}`
    let memberFooterElement
    if (isMember) {
        memberFooterElement = (
            <p className="card-footer-item has-background-success-light">Membership Verified</p>
        )
    } else {
        memberFooterElement = (
            <p className="card-footer-item">Non-member</p>
        )
    }
    return (
        <div className="column is-one-quarter">
            <div className="card channel">
                <div className="card-image">
                    <figure className="channel-thumbnail image is-128x128">
                        <a href={channelHref} target="_blank" rel="noopener noreferrer">
                            <img src={doc.Thumbnail} alt="channel thumbnail" className="is-rounded mt-2" />
                        </a>
                    </figure>
                </div>
                <div className="card-content has-text-centered">
                    <div className="content">
                        <h4 className="title is-6">{doc!.ChannelTitle}</h4>
                    </div>
                </div>
                <footer v-else className="card-footer">
                    {memberFooterElement}
                </footer>
            </div>
        </div>
    )
}

export default YouTubeMembershipCard