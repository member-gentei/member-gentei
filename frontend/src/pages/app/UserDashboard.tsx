import { Fragment } from "react";
import AdminServers from "../../components/dashboard/AdminServers";
import DiscordServers from "../../components/dashboard/DiscordServers";
import SelfManage from "../../components/dashboard/SelfManage";
import { YouTubeLoginOverlay } from "../../components/dashboard/YouTubeLogin";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";
import { Grid } from "@mui/joy";

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
  if (store.user === undefined || !store.user.YouTube.Valid) {
    serversOrLoginNode = (
      <section>
        <DiscordServers />
        <YouTubeLoginOverlay />
      </section>
    );
  } else {
    serversOrLoginNode = (
      <section>
        <DiscordServers />
      </section>
    );
  }
  return (
    <Grid container rowSpacing={1}>
      <Grid xs={12}>
        {serversOrLoginNode}
        {loginRequiredOverlay}
      </Grid>
      <Grid xs={12}>
        <SelfManage />
      </Grid>
      <Grid xs={12}>
        <AdminServers />
      </Grid>
    </Grid>
  );
}
