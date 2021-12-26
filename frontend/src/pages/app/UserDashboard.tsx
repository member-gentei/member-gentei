import { Fragment } from "react";
import AdminServers from "../../components/dashboard/AdminServers";
import DiscordServers from "../../components/dashboard/DiscordServers";

export default function UserDashboard() {
  return (
    <Fragment>
      <DiscordServers />
      <AdminServers />
    </Fragment>
  );
}
