import { Fragment } from "react";
import { RiCheckFill, RiCloseFill } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { LoadState, ZeroTime } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import { Talent, useTalents } from "../../stores/TalentStore";
import { useUser } from "../../stores/UserStore";
import DiscordServerImg from "../DiscordServerImg";
import {
  Alert,
  AspectRatio,
  Card,
  CardContent,
  Divider,
  Grid,
  Link,
  Skeleton,
  Typography,
} from "@mui/joy";
import WarningIcon from "@mui/icons-material/Warning";

export default function DiscordServers() {
  const [userStore] = useUser();
  let serverColumns;
  let uncheckedNotice = null;
  if (userStore.userLoad <= LoadState.Started) {
    serverColumns = (
      <Skeleton>
        <Grid container>
          <Grid xs={6}></Grid>
          <Grid xs={6}></Grid>
        </Grid>
      </Skeleton>
    );
  } else if (userStore.derived.sortedServers.length > 0) {
    serverColumns = (
      <Grid container sx={{ mt: 2, mb: 2 }} spacing={1}>
        {userStore.derived.sortedServers.map((serverID) => (
          <Grid xs={6} key={`dsr-${serverID}`}>
            <DiscordServerWithRoles id={serverID} />
          </Grid>
        ))}
      </Grid>
    );
  }
  if (
    userStore.user?.LastRefreshed === ZeroTime &&
    !!userStore.user?.YouTube.ID
  ) {
    uncheckedNotice = (
      <Alert startDecorator={<WarningIcon />} color="warning">
        <div>
          <div>Membership check not finished</div>
          <Typography level="body-sm" color="warning">
            The role assignments below do not yet reflect your current YouTube
            memberships. The job scheduled to check your memberships has not
            finished - this message will disappear after it has.
          </Typography>
        </div>
      </Alert>
    );
  }
  return (
    <Fragment>
      <Typography level="h2">Servers and Roles</Typography>
      <Typography>
        Servers that you've joined that participate in Gentei's members-only
        role management are listed below.
      </Typography>
      <p className="mb-4"></p>
      {uncheckedNotice}
      {serverColumns}
      <Typography>
        If a server you've joined is not shown above <strong>and</strong>{" "}
        <code>/gentei</code> is a slash command on that server, please wait a
        few minutes for the bot to refresh server memberships. Discord can take
        a few minutes to make server information available to integrations like
        Gentei.
      </Typography>
    </Fragment>
  );
}

interface DiscordServerRoleProps {
  id: string;
}

export function DiscordServerWithRoles(props: DiscordServerRoleProps) {
  return (
    <GuildContainer isGlobal scope={props.id}>
      <DiscordServerWithRolesInner {...props} />
    </GuildContainer>
  );
}

function DiscordServerWithRolesInner({ id }: DiscordServerRoleProps) {
  const [userStore] = useUser();
  const [talentStore, talentActions] = useTalents();
  const [guildStore, actions] = useGuild();
  actions.load(id);
  talentActions.loadAll();
  if (guildStore.guildState <= LoadState.Started) {
    return <span className="spinner mx-auto"></span>;
  }
  const guild = guildStore.guild!;
  const serverURL = `https://discord.com/channels/${id}`;
  let membershipNode;
  let memberships = Object.entries(guildStore.guild?.RolesByTalent || {}).map(
    ([k, v]) => {
      const talentID = k;
      const meta = (userStore.user?.Memberships || {})[k];
      return (
        <Grid>
          <RoleMembership
            key={`${id}-${talentID}`}
            talent={talentStore.talentsByID[talentID]}
            roleName={v!.Name}
            verifyTime={meta?.Failed ? 0 : meta?.LastVerified}
          />
        </Grid>
      );
    }
  );
  if (memberships.length === 0) {
    membershipNode = (
      <Typography>
        This server has not configured memberships yet. Please be discreet until
        server moderation announces something!
      </Typography>
    );
  } else {
    membershipNode = (
      <div>
        {/* <Typography sx={{ fontWeight: 600 }}>Discord roles</Typography> */}
        <Grid container spacing={1}>
          {memberships}
        </Grid>
      </div>
    );
  }
  let iconNode;
  if (guild.Icon.length > 0) {
    iconNode = (
      <DiscordServerImg
        guildID={guild.ID}
        imgHash={guild.Icon}
        className="is-rounded"
      />
    );
  } else {
    iconNode = <SiDiscord size={48} />;
  }
  return (
    <Card orientation="horizontal">
      {iconNode}
      <CardContent>
        <Link href={serverURL}>{guild.Name}</Link>
        {membershipNode}
      </CardContent>
    </Card>
  );
}

interface RoleMembershipProps {
  talent?: Talent;
  roleName: string;
  verifyTime?: number;
}

function RoleMembership({ talent, roleName, verifyTime }: RoleMembershipProps) {
  if (talent === undefined) {
    return (
      <Skeleton>
        <Card>
          <Skeleton>lorem ipsum lmao</Skeleton>
        </Card>
      </Skeleton>
    );
  }
  const channelURL = `https://youtube.com/channel/${talent.ID}`;
  let endDecorator = null;
  let tooltip = null;
  if (!!verifyTime) {
    const verifyTs = new Date(verifyTime! * 1000);
    const verifyTimeStr = `${verifyTs.toDateString()} ${verifyTs.toTimeString()}`;
    tooltip = `Last verified at ${verifyTimeStr}`;
    endDecorator = <RiCheckFill color="green" />;
  } else {
    endDecorator = <RiCloseFill color="red" />;
  }
  return (
    <Card sx={{ alignItems: "center", textAlign: "center" }}>
      <a href={channelURL} title={`YouTube channel for ${talent.Name}`}>
        <AspectRatio ratio="1" sx={{ width: 128 }}>
          <img
            className="rounded"
            src={talent.Thumbnail}
            alt={`Channel icon for ${talent.Name}`}
          />
        </AspectRatio>
      </a>
      <Divider />
      <Typography endDecorator={endDecorator}>
        <span className="discord-mention">@{roleName}</span>
      </Typography>
    </Card>
  );
}
