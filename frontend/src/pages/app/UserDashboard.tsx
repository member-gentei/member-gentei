import { Fragment } from "react";
import AdminServers from "../../components/dashboard/AdminServers";
import DiscordServers from "../../components/dashboard/DiscordServers";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";

export default function UserDashboard() {
  const [store] = useUser();
  let loginRequiredOverlay = null;
  if (store.user === undefined && store.userLoad > LoadState.NotStarted) {
    loginRequiredOverlay = (
      <div className="overlay is-flex is-justify-content-center is-align-content-center is-align-items-center">
        <div className="message is-info">
          <div className="message-body">
            Please sign in or register to view and manage memberships.
          </div>
        </div>
      </div>
    );
  }
  return (
    <Fragment>
      <div style={{ position: "relative" }}>
        <DiscordServers />
        <AdminServers />
        {loginRequiredOverlay}
      </div>
    </Fragment>
  );
}
