import firebase from "firebase/app"
import "firebase/firestore"
import { useState } from "react"

interface ChannelConnectionProps {
    user: firebase.User
    channelID?: string
}

const ChannelConnection = ({ user, channelID }: ChannelConnectionProps) => {
    const youTubeLoginURL = "https://accounts.google.com/o/oauth2/v2/auth?client_id=649732146530-s4cj4tqo2impojg7ljol2chsuj1us81s.apps.googleusercontent.com&redirect_uri=https%3A%2F%2Fmember-gentei.tindabox.net%2Flogin%2Fyoutube&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fyoutube.force-ssl&access_type=offline&prompt=consent"
    const [disconnecting, setDisconnecting] = useState(false)
    if (!!channelID) {
        const channelHref = `https://youtube.com/channel/${channelID}`
        const dcYouTubeAccount = () => {
            firebase.firestore().collection("users").doc(user.uid)
                .collection("private").doc("youtube").delete()
            setDisconnecting(true)
        }
        let buttonClassName = "button is-danger"
        if (disconnecting) {
            buttonClassName += " is-loading"
        }
        return (
            <div className="content">
                <p>Connected channel: <a href={channelHref} target="_blank" rel="noreferrer">{channelHref}</a></p>
                <div className="buttons is-centered">
                    <button
                        className={buttonClassName}
                        onClick={dcYouTubeAccount}>Disconnect Channel</button>
                </div>
            </div>
        )
    }
    return (
        <div className="content">
            <p>Please connect your YouTube account to verify your membership(s).</p>
            <div className="has-text-centered">
                <a href={youTubeLoginURL}><div className="signin-google-button"></div></a>
            </div>
        </div>
    )
}

export default ChannelConnection