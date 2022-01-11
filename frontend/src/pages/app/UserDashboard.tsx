import { Fragment } from "react";
import AdminServers from "../../components/dashboard/AdminServers";
import DiscordServers from "../../components/dashboard/DiscordServers";
import YouTubeLogin from "../../components/dashboard/YouTubeLogin";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";

export default function UserDashboard() {
  const [store] = useUser();
  let loginRequiredOverlay = null;
  if (store.user === undefined && store.userLoad > LoadState.Started) {
    loginRequiredOverlay = (
      <div className="overlay is-flex is-justify-content-center is-align-content-center is-align-items-center">
        <div className="message is-info">
          <div className="message-body">
            <div className="content">
              <p>Please sign in or register to view and manage memberships.</p>
              <p>
                To verify your membership(s), you will have to connect your
                YouTube account after registering.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }
  let serversOrLoginNode;
  if (store.user === undefined || store.user.YouTube.Valid) {
    serversOrLoginNode = (
      <section className="section">
        <DiscordServers />
      </section>
    );
  } else {
    serversOrLoginNode = (
      <section className="section" style={{ position: "relative" }}>
        <DiscordServers />
        <YouTubeLogin />
      </section>
    );
  }
  return (
    <Fragment>
      <div style={{ position: "relative" }}>
        {serversOrLoginNode}
        <AdminServers />
        {loginRequiredOverlay}
      </div>
    </Fragment>
  );
}
