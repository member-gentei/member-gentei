import React, { ReactNode } from "react";
import { RiCheckFill, RiCloseFill } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { LoadState } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import { Talent, useTalents } from "../../stores/TalentStore";
import { useUser } from "../../stores/UserStore";
import styles from "./DiscordServers.module.css";

export default function DiscordServers() {
  const [userStore] = useUser();
  let serverColumns;
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
  return (
    <section className="section">
      <div className="container">
        <h1 className="title is-2">Servers and Roles</h1>
        <p className="mb-4">
          Servers that you've joined that have members-only role management are
          listed below.
        </p>
        {serverColumns}
        <p>
          If a server you've joined is not shown above <strong>and</strong>{" "}
          <code>/gentei</code> is a slash command on that server, please wait a
          few minutes for the bot to refresh server memberships. Discord can
          take a few minutes to make information available to integrations.
        </p>
      </div>
    </section>
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
  // const memberships = guildStore.guild?.Settings?.RoleMapping.map((m) => (
  //   <RoleMembership key={`${id}-${m.RoleName}`} {...m} />
  // ));
  let memberships = Object.entries(
    guildStore.guild?.Settings?.RoleMapping || {}
  ).map(([k, v]) => (
    <RoleMembership
      key={`${k}-${v!.ID}`}
      talent={talentStore.talentsByID[v!.ID]!}
      roleName={v!.Name}
      verifyTime={(userStore.user?.Roles || {})[v!.ID]}
    />
  ));
  if (memberships.length === 0) {
    memberships = [
      <p>
        This server has not configured memberships yet. Please be discreet until
        server moderation announces something!
      </p>,
    ];
  }
  let iconNode;
  if (guild.Icon.length > 0) {
    const iconURL = `https://cdn.discordapp.com/icons/${id}/${guild.Icon}.webp?size=128`;
    // bonus: make 'em gifs if applicable
    let onHover: React.MouseEventHandler<HTMLImageElement> = () => {};
    let offHover: React.MouseEventHandler<HTMLImageElement> = () => {};
    if (guild.Icon.startsWith("a_")) {
      const gifURL = iconURL.replace(".webp", ".gif");
      onHover = (e) => {
        e.currentTarget.setAttribute("src", gifURL);
      };
      offHover = (e) => {
        e.currentTarget.setAttribute("src", iconURL);
      };
    }
    iconNode = (
      <img
        className="is-rounded"
        src={iconURL}
        alt="Discord server icon"
        onMouseOver={onHover}
        onMouseOut={offHover}
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
              <a href={serverURL}>{iconNode}</a>
            </p>
          </figure>
          <div className="media-content">
            <div className="content">
              <p>
                <a href={serverURL} title="Link to Discord server">
                  <strong>{guild.Name}</strong>
                </a>
              </p>
            </div>
            {memberships}
          </div>
        </div>
      </div>
    </div>
  );
}

interface RoleMembershipProps {
  talent: Talent;
  roleName: string;
  verifyTime?: number;
}

function RoleMembership({ talent, roleName, verifyTime }: RoleMembershipProps) {
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
            <RiCheckFill />
          </span>
        </span>
      </div>
    );
  } else {
    footerItem = (
      <div className="card-footer-item">
        <span className="icon-text">
          <span>@{roleName}</span>
          <span className="icon">
            <RiCloseFill />
          </span>
        </span>
      </div>
    );
  }
  return (
    <div className={`card m-2 ${styles.verifyCard}`}>
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
