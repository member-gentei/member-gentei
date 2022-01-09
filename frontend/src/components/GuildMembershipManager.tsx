import { Fragment } from "react";
import { RiCheckboxCircleFill, RiCloseCircleLine } from "react-icons/ri";
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
  let roleNode;
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
    roleNode = <code>-</code>;
  } else {
    statusNode = (
      <Fragment>
        <span className="icon-text">
          <span className="icon">
            <RiCheckboxCircleFill color="green" size={24} />
          </span>
          <span>Role management in effect</span>
        </span>
      </Fragment>
    );
    roleNode = (
      <Fragment>
        <span className="discord-role">@{roleMapping.Name}</span> (
        <code>{roleMapping.ID}</code>)
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
                <th>Discord Role</th>
                <td>
                  {roleNode}
                  <span className="help">
                    To change, use{" "}
                    <code>/gentei manage map {talentID} @newrole</code>.
                  </span>
                </td>
              </tr>
              <tr>
                <th>Membership count</th>
                <td>
                  {roleMapping ? <span>? members</span> : <code>-</code>}
                  <span className="help">
                    Members in this server with the role. Count refreshed daily.
                  </span>
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
