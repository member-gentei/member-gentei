import firebase from 'firebase/app';
import "firebase/auth"
import { useEffect } from 'react';
import { useHistory, useLocation, useParams } from "react-router-dom"

interface LoginProps {
    onDiscordLogin: (status: LoginStatus) => void
}

export interface LoginStatus {
    status: string
    ok: boolean
}

const Login = ({ onDiscordLogin }: LoginProps) => {
    const { loginType } = useParams<{ loginType: string }>()
    const history = useHistory()
    const code = new URLSearchParams(useLocation().search).get("code")
    if (!code) {
        console.log("no OAuth code found")
        history.push("/app")
    }
    useEffect(() => {
        if (loginType === "discord") {
            handleDiscordLogin(onDiscordLogin, code!, history)
        }
        // we don't want to do this more than once.
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])
    return (
        <div className="login-spinner has-text-centered">
            <div className="lds-ring"><div></div><div></div><div></div><div></div></div>
            <p>Logging you in...</p>
        </div>
    )
}

async function handleDiscordLogin(
    onDiscordLogin: (status: LoginStatus) => void,
    code: string,
    history: ReturnType<typeof useHistory>
) {
    const loginURL = "https://us-central1-member-gentei.cloudfunctions.net/Auth?service=discord"
    let r: Response
    try {
        let params = new URLSearchParams()
        params.set("code", code)
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
            const body: { jwt: string } = await r.json()
            await firebase.auth().signInWithCustomToken(body["jwt"])
            onDiscordLogin({ status: "", ok: true })
        } else {
            onDiscordLogin({ status: await r.text(), ok: false })
        }
    } catch (e) {
        console.error(e)
        onDiscordLogin({ status: "Unexpected error signing in - probably a browser CORS error?", ok: false })
    }
    history.push("/app")
}

export default Login