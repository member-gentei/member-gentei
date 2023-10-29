import { SiDiscord } from "react-icons/si";
import { Link } from "react-router-dom";
import { LoadState } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import { useUser } from "../../stores/UserStore";
import DiscordServerImg from "../DiscordServerImg";
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  Grid,
  Stack,
  Table,
  Typography,
} from "@mui/joy";

export default function AdminServers() {
  const [store, actions] = useUser();
  actions.getMe();
  if (store.userLoad <= LoadState.Started) {
    return (
      <div className="section">
        <div className="container">
          <Typography level="h2">Administration</Typography>
          <div className="has-text-centered`">
            <span className="spinner mx-auto"></span>
          </div>
        </div>
      </div>
    );
  }
  let serverGrid;
  if (!!store.user?.ServerAdmin) {
    serverGrid = (
      <Grid container spacing={1}>
        {store.user.ServerAdmin.map((guildID) => (
          <Grid xs={12} md={6} key={`${guildID}-admin`}>
            <AdminServerSnippet guildID={guildID} />
          </Grid>
        ))}
      </Grid>
    );
  } else {
    serverGrid = (
      <p className="content">
        You do not have permission to configure or audit Gentei for any Discord
        server integrations.
      </p>
    );
  }
  return (
    <Stack component="section" spacing={2} mb={2}>
      <Typography level="h2">Administration</Typography>
      {serverGrid}
      <Box sx={{ textAlign: "center" }}>
        <Link to="enroll" className="button spin-hover">
          <Button size="lg" endDecorator={<SiDiscord className="spin-me" />}>
            Enroll a new Discord server
          </Button>
        </Link>
      </Box>
    </Stack>
  );
}

interface AdminServerSnippetProps {
  guildID: string;
}

function AdminServerSnippet(props: AdminServerSnippetProps) {
  return (
    <GuildContainer isGlobal scope={props.guildID}>
      <AdminServerSnippetInner {...props} />
    </GuildContainer>
  );
}

function AdminServerSnippetInner({ guildID }: AdminServerSnippetProps) {
  const [store, actions] = useGuild();
  actions.load(guildID);
  if (store.guildState <= LoadState.Started) {
    return (
      <div className="box">
        <div className="has-text-centered">
          <span className="spinner mx-auto"></span>
        </div>
      </div>
    );
  }
  const guild = store.guild!;
  let iconNode;
  if (guild.Icon.length > 0) {
    iconNode = (
      <DiscordServerImg
        guildID={guild.ID}
        imgHash={guild.Icon}
        size={128}
        className="is-rounded"
      />
    );
  } else {
    iconNode = <SiDiscord size={128} />;
  }
  return (
    <Card>
      <CardContent orientation="horizontal">
        {iconNode}
        <Table borderAxis="bothBetween">
          <tbody>
            <tr>
              <th style={{ width: "10rem" }}>Name</th>
              <th>{guild.Name}</th>
            </tr>
            <tr>
              <th>ID / Link</th>
              <td>
                <a href={`https://discord.com/channels/${guild.ID}`}>
                  {guild.ID}
                </a>
              </td>
            </tr>
            <tr>
              <th>Membership roles</th>
              <td>{Object.keys(guild.RolesByTalent).length}</td>
            </tr>
          </tbody>
        </Table>
      </CardContent>
      <CardActions buttonFlex="0">
        <Button
          component="a"
          href={`server/${guild.ID}`}
          size="sm"
          variant="soft"
        >
          View/Edit
        </Button>
      </CardActions>
    </Card>
  );
}
