import firebase from 'firebase/app';
import "firebase/auth"
import { useEffect } from 'react';
import { useHistory, useLocation, useParams } from "react-router-dom"

interface LoginProps {
    ready: boolean
    onDiscordLogin: (status: LoginStatus) => void
    onYouTubeLogin: (status: LoginStatus) => void
}

export interface LoginStatus {
    status: string
    ok: boolean
}

const Login = ({ ready, onDiscordLogin, onYouTubeLogin }: LoginProps) => {
    const { loginType } = useParams<{ loginType: string }>()
    const history = useHistory()
    const code = new URLSearchParams(useLocation().search).get("code")
    if (!code) {
        console.log("no OAuth code found")
        history.push("/app")
    }
    let loadingMessage = "Logging you in..."
    useEffect(() => {
        if (!ready) {
            loadingMessage = "Loading..."
            return
        }
        if (loginType === "discord" || loginType === "youtube") {
            handleLogin(onDiscordLogin, onYouTubeLogin, loginType, code!, history)
        }
        // we don't want to do this dance more than once.
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [ready])
    return (
        <div className="login-spinner has-text-centered">
            <div className="lds-ring"><div></div><div></div><div></div><div></div></div>
            <p>{loadingMessage}</p>
        </div>
    )
}

async function handleLogin(
    onDiscordLogin: (status: LoginStatus) => void,
    onYouTubeLogin: (status: LoginStatus) => void,
    loginType: string,
    code: string,
    history: ReturnType<typeof useHistory>
) {
    const loginURL = "https://us-central1-member-gentei.cloudfunctions.net/Auth?service=" + loginType
    let r: Response
    try {
        let params = new URLSearchParams()
        params.set("code", code)
        if (loginType === "youtube") {
            const idToken = await getFirebaseUserToken()
            if (typeof idToken == "undefined") {
                onYouTubeLogin({ status: "unable to refresh current credentials", ok: false })
                history.push("/app")
                return
            }
            params.set("jwt", idToken!)
        }
        r = await fetch(loginURL, {
            method: "POST",
            headers: {
                "X-Requested-With": "XMLHttpRequest",
                "Content-Type": "application/x-www-form-urlencoded",
                "Accept": "application/json"
            },
            body: params,
        })
        if (r.ok) {
            if (loginType === "youtube") {
                onYouTubeLogin({ status: "", ok: true })
            } else {
                const body: { jwt: string } = await r.json()
                await firebase.auth().signInWithCustomToken(body["jwt"])
                onDiscordLogin({ status: "", ok: true })
            }
        } else {
            if (loginType === "youtube") {
                onYouTubeLogin({ status: await r.text(), ok: false })
            } else {
                onDiscordLogin({ status: await r.text(), ok: false })
            }
        }
    } catch (e) {
        console.error(e)
        if (loginType === "youtube") {
            onYouTubeLogin({ status: "Unexpected error signing in to YouTube.", ok: false })
        } else {
            onDiscordLogin({ status: "Unexpected error signing in - probably a browser CORS error?", ok: false })
        }
    }
    history.push("/app")
}

const getFirebaseUserToken = async () => {
    const currentUser = firebase.auth().currentUser
    if (currentUser !== null) {
        return await currentUser.getIdToken()
    }
    return await new Promise<string>((resolve, reject) => {
        firebase.auth().onAuthStateChanged(function (user) {
            if (user == null) {
                return
            }
            user.getIdToken().then(resolve)
        })
    })
}

export default Login