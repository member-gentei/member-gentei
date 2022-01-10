import React, { ReactNode } from "react";
import { RiCheckFill, RiCloseFill } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { LoadState, ZeroTime } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import { Talent, useTalents } from "../../stores/TalentStore";
import { useUser } from "../../stores/UserStore";
import DiscordServerImg from "../DiscordServerImg";

export default function DiscordServers() {
  const [userStore] = useUser();
  let serverColumns;
  let uncheckedNotice = null;
  if (userStore.userLoad <= LoadState.Started) {
    serverColumns = (
      <div className="columns is-multiline">
        <div className="column has-text-centered">
          <span className="spinner mx-auto"></span>
        </div>
      </div>
    );
  } else if (userStore.derived.sortedServers.length > 0) {
    serverColumns = (
      <div className="columns is-multiline">
        {userStore.derived.sortedServers.map((serverID) => (
          <div key={`dsr-${serverID}`} className="column is-half">
            <DiscordServerWithRoles id={serverID} />
          </div>
        ))}
      </div>
    );
  }
  if (userStore.user?.LastRefreshed === ZeroTime) {
    uncheckedNotice = (
      <div className="columns is-centered">
        <div className="column is-half">
          <div className="message is-warning">
            <div className="message-header">
              <p>Membership check not finished</p>
            </div>
            <div className="message-body">
              <p>
                The role assignments below do not yet reflect your current
                YouTube memberships. The job scheduled to check your memberships
                has not finished - this message will disappear after it has.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="container">
      <h1 className="title is-2">Servers and Roles</h1>
      <p className="mb-4">
        Servers that you've joined that have members-only role management are
        listed below.
      </p>
      {uncheckedNotice}
      {serverColumns}
      <p>
        If a server you've joined is not shown above <strong>and</strong>{" "}
        <code>/gentei</code> is a slash command on that server, please wait a
        few minutes for the bot to refresh server memberships. Discord can take
        a few minutes to make membership information available to integrations
        like Gentei.
      </p>
    </div>
  );
}

interface DiscordServerRoleProps {
  id: string;
}

function DiscordServerWithRoles(props: DiscordServerRoleProps) {
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
  let memberships = Object.entries(
    guildStore.guild?.Settings?.RoleMapping || {}
  ).map(([k, v]) => {
    const talentID = k;
    return (
      <RoleMembership
        key={`${id}-${talentID}`}
        talent={talentStore.talentsByID[talentID]}
        roleName={v!.Name}
        verifyTime={(userStore.user?.Roles || {})[v!.ID]}
      />
    );
  });
  if (memberships.length === 0) {
    membershipNode = (
      <p className="content">
        This server has not configured memberships yet. Please be discreet until
        server moderation announces something!
      </p>
    );
  } else {
    membershipNode = (
      <div className="content">
        <span className="is-size-6 has-text-weight-bold">Discord roles</span>
        <div className="is-flex is-flex-wrap-wrap">{memberships}</div>
      </div>
    );
  }
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
    iconNode = <SiDiscord size={64} />;
  }
  return (
    <div className="card">
      <div className="card-content">
        <div className="media">
          <figure className="media-left">
            <p className="image is-64x64">
              <a href={serverURL} title="Link to Discord server">
                {iconNode}
              </a>
            </p>
          </figure>
          <div className="media-content">
            <div className="content">
              <p>
                <a
                  className="is-size-5 has-text-weight-bold"
                  href={serverURL}
                  title="Link to Discord server"
                >
                  {guild.Name}
                </a>
              </p>
            </div>
            {membershipNode}
          </div>
        </div>
      </div>
    </div>
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
      <div className="card m-2">
        <div className="card-content has-text-centered">
          <span className="spinner mx-auto"></span>
        </div>
      </div>
    );
  }
  const channelURL = `https://youtube.com/channel/${talent.ID}`;
  let footerItem: ReactNode;
  if (!!verifyTime) {
    const verifyTs = new Date(verifyTime! * 1000);
    const verifyTimeStr = `${verifyTs.toDateString()} ${verifyTs.toTimeString()}`;
    const tooltip = `Last verified at ${verifyTimeStr}`;
    footerItem = (
      <div className="card-footer-item has-background-success-light">
        <span className="icon-text">
          <span>@{roleName}</span>
          <span
            className="icon has-tooltip-arrow has-text-success-dark"
            data-tooltip={tooltip}
          >
            <RiCheckFill color="green" />
          </span>
        </span>
      </div>
    );
  } else {
    footerItem = (
      <div className="card-footer-item">
        <span className="icon-text">
          <span className="discord-role">@{roleName}</span>
          <span className="icon">
            <RiCloseFill color="red" />
          </span>
        </span>
      </div>
    );
  }
  return (
    <div className="card m-2">
      <div className="card-image">
        <a href={channelURL} title={`YouTube channel for ${talent.Name}`}>
          <figure className="image is-128x128">
            <img className="is-rounded" src={talent.Thumbnail} alt="" />
          </figure>
        </a>
      </div>
      <div className="card-footer">{footerItem}</div>
    </div>
  );
}
