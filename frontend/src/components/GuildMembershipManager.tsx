import { Fragment } from "react";
import { RiCheckboxCircleFill, RiCloseCircleLine } from "react-icons/ri";
import { useGuild } from "../stores/GuildStore";
import styles from "./GuildMembershipManager.module.css";
import TalentCard from "./TalentCard";
import { Card, Grid, Table, Typography } from "@mui/joy";

interface GuildMembershipManagerProps {
  talentID: string;
  settings?: {};
}

export default function GuildMembershipManager({
  talentID,
}: GuildMembershipManagerProps) {
  const [state] = useGuild();
  const guild = state.guild!;
  const roleMapping = (guild.RolesByTalent || {})[talentID];
  let statusNode;
  let roleNode;
  if (roleMapping === undefined) {
    statusNode = (
      <Fragment>
        <Typography
          startDecorator={<RiCloseCircleLine color="red" size={24} />}
        >
          Not yet configured <br />
        </Typography>
        <Typography level="body-xs">
          Please run{" "}
          <code>
            /gentei manage map youtube-channel-id:{talentID} role:@role
          </code>
        </Typography>
      </Fragment>
    );
    roleNode = <code>-</code>;
  } else {
    statusNode = (
      <Typography
        startDecorator={<RiCheckboxCircleFill color="green" size={24} />}
      >
        Exclusive role management in effect
      </Typography>
    );
    roleNode = (
      <Fragment>
        <span className="discord-role">@{roleMapping.Name}</span> (
        <code>{roleMapping.ID}</code>)
      </Fragment>
    );
  }
  return (
    <Grid container spacing={4}>
      <Grid xs={2}>
        <TalentCard cardClassNames={[styles.height200]} channelID={talentID} />
      </Grid>
      <Grid xs>
        <Card>
          <Table
            sx={{
              "& th": { width: "10rem" },
              overflow: "hidden",
            }}
          >
            <tbody>
              <tr>
                <th>Status</th>
                <td>{statusNode}</td>
              </tr>
              <tr>
                <th>Discord Role</th>
                <td>
                  {roleNode}
                  <Typography level="body-xs">
                    To change, use{" "}
                    <code>
                      /gentei manage map youtube-channel-id:{talentID}{" "}
                      role:@role
                    </code>
                    .
                  </Typography>
                </td>
              </tr>
              <tr>
                <th>Membership count</th>
                <td>
                  {roleMapping ? <span>? members</span> : <code>-</code>}
                  <Typography level="body-xs">
                    Members in this server with the role. Count refreshed daily.
                  </Typography>
                </td>
              </tr>
            </tbody>
          </Table>
          <div className="message is-info">
            <div className="message-body">Customization coming soon!</div>
          </div>
        </Card>
      </Grid>
    </Grid>
  );
}
