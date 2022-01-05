import { Fragment } from "react";
import { RiCloseCircleLine } from "react-icons/ri";
import { useGuild } from "../stores/GuildStore";
import styles from "./GuildMembershipManager.module.css";
import TalentCard from "./TalentCard";

interface GuildMembershipManagerProps {
  talentID: string;
  settings?: {};
}

export default function GuildMembershipManager({
  talentID,
}: GuildMembershipManagerProps) {
  const [state] = useGuild();
  const guild = state.guild!;
  const roleMapping = (guild.Settings?.RoleMapping || {})[talentID];
  let statusNode;
  if (roleMapping === undefined) {
    statusNode = (
      <Fragment>
        <span className="icon-text">
          <span className="icon">
            <RiCloseCircleLine color="red" size={24} />
          </span>
          <span>
            Not yet configured <br />
            Please run{" "}
            <code>/gentei manage map UCAL_ZudIZXhCDrniD4ZQobQ @role</code>
          </span>
        </span>
      </Fragment>
    );
  }
  return (
    <div className="columns is-gapless">
      <div className="column is-narrow">
        <TalentCard cardClassNames={[styles.height200]} channelID={talentID} />
      </div>
      <div className="column">
        <div className={`box m-1 ${styles.scrolly}`}>
          <table className="table is-narrow" style={{ width: "auto" }}>
            <tbody>
              <tr>
                <th>Status</th>
                <td>{statusNode}</td>
              </tr>
              <tr>
                <th>Discord Role ID</th>
                <td>
                  <code>{roleMapping || "-"}</code>
                  <span className="help">
                    To see the role name, use <code>/gentei info</code>.
                  </span>
                </td>
              </tr>
              <tr>
                <th>Membership count</th>
                <td>
                  {roleMapping ? (
                    <span>{roleMapping} Members</span>
                  ) : (
                    <code>-</code>
                  )}
                  <span className="help">Refreshed daily</span>
                </td>
              </tr>
            </tbody>
          </table>
          <div className="message is-info">
            <div className="message-body">Customization coming soon!</div>
          </div>
        </div>
      </div>
    </div>
  );
}
