import { Fragment } from "react";
import { Link as RRDLink, useParams } from "react-router-dom";
import GuildMembershipManager from "../../components/GuildMembershipManager";
import { LoadState } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import {
  Box,
  Breadcrumbs,
  Button,
  Grid,
  Link,
  Stack,
  Table,
  Typography,
} from "@mui/joy";

export default function GuildAdmin() {
  const { guildID } = useParams();
  if (!guildID) {
    return (
      <div className="has-text-centered">
        <span className="spinner"></span>
      </div>
    );
  }
  return (
    <GuildContainer isGlobal scope={guildID}>
      <GuildAdminInner guildID={guildID} />
    </GuildContainer>
  );
}

interface GuildAdminInnerProps {
  guildID: string;
}

function GuildAdminInner({ guildID }: GuildAdminInnerProps) {
  const [store, actions] = useGuild();
  actions.load(guildID);
  if (store.guildState <= LoadState.Started) {
    return (
      <div className="has-text-centered">
        <span className="spinner"></span>
      </div>
    );
  }
  const guild = store.guild!;
  const guildURL = `https://discord.com/channels/${guildID}`;
  const membershipManagers = (guild.TalentIDs || []).map((talentID) => (
    <GuildMembershipManager key={`manage-${talentID}`} talentID={talentID} />
  ));
  return (
    <Stack component="section" spacing={2} sx={{ mb: 2 }}>
      <Breadcrumbs>
        <RRDLink to="/app">
          <Link>Home</Link>
        </RRDLink>
        <Typography>{guild.Name}</Typography>
      </Breadcrumbs>
      <Typography level="h2">{guild.Name}</Typography>
      <Table>
        <tbody>
          <tr>
            <th>Server ID and Link</th>
            <td>
              <Link
                href={guildURL}
                target="_blank"
                rel="noreferrer"
                title="Open server in a new tab"
              >
                {guild.ID}
              </Link>
            </td>
          </tr>
          <tr>
            <th>Server Name</th>
            <td>{guild.Name}</td>
          </tr>
          <tr>
            <th>Discord IDs with admin access</th>
            <td>
              {guild.AdminIDs.map((snowflake) => (
                <Fragment key={snowflake}>
                  <code>{snowflake}</code>{" "}
                </Fragment>
              ))}
            </td>
          </tr>
          <tr>
            <th>Audit log channel ID and link</th>
            <td>
              {guild.AuditLogChannelID ? (
                <a
                  href={`${guildURL}/${guild.AuditLogChannelID}`}
                  target="_blank"
                  rel="noreferrer"
                  title="Open channel in a new tab"
                >
                  {guild.AuditLogChannelID}
                </a>
              ) : (
                <code>- (none)</code>
              )}
              <Typography level="body-sm">
                To change, use{" "}
                <code>/gentei manage audit-log channel:#channel</code>
              </Typography>
            </td>
          </tr>
        </tbody>
      </Table>
      <Typography>
        For other information and nicely formatted{" "}
        <span className="discord-role">@mentions</span>, please run{" "}
        <code>/gentei info</code> in your server.
      </Typography>
      <Box>
        <Typography level="h3">Memberships</Typography>
        <Typography>
          Settings that are hard to configure using slash commands can be edited
          below.
        </Typography>
        <Stack spacing={2}>{membershipManagers}</Stack>
      </Box>
      <Box sx={{ textAlign: "center" }}>
        <RRDLink
          to={{
            pathname: "/app/enroll",
            search: new URLSearchParams({ server: guildID }).toString(),
          }}
        >
          <Button>Add/Remove Memberships</Button>
        </RRDLink>
      </Box>
    </Stack>
  );
}
