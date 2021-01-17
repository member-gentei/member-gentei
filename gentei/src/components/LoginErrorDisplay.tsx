import { LoginStatus } from "../pages/Login"


const LoginErrorDisplay = ({ loginStatus }: { loginStatus?: LoginStatus }) => {
    if (typeof loginStatus != "undefined" && !loginStatus.ok) {
        return (
            <div className="has-text-danger-dark has-text-weight-medium">⚠️ { loginStatus!.status}</div>
        )
    }
    return null

}

export default LoginErrorDisplay