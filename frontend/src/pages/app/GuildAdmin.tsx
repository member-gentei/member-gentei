import { Fragment } from "react";
import { Link, useParams } from "react-router-dom";
import GuildMembershipManager from "../../components/GuildMembershipManager";
import { LoadState } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";

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
    <section className="section">
      <div className="container">
        <nav className="breadcrumb">
          <ul>
            <li>
              <Link to="/app">Home</Link>
            </li>
            <li className="is-active">
              <Link to="#">{guild.Name}</Link>
            </li>
          </ul>
        </nav>
        <h1 className="title">{guild.Name}</h1>
        <div className="content">
          <dl>
            <dt>Server ID and Link</dt>
            <dd>
              <a
                href={guildURL}
                target="_blank"
                rel="noreferrer"
                title="Open server in a new tab"
              >
                {guild.ID}
              </a>
            </dd>
            <dt>Server Name</dt>
            <dd>{guild.Name}</dd>
            <dt>Discord IDs with admin access</dt>
            <dd>
              {guild.AdminIDs.map((snowflake) => (
                <Fragment key={snowflake}>
                  <code>{snowflake}</code>{" "}
                </Fragment>
              ))}
            </dd>
            <dt>Audit log channel ID and link</dt>
            <dd>
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

              <p className="help">
                To change, use{" "}
                <code>/gentei manage audit-log channel:#channel</code>
              </p>
            </dd>
          </dl>
          <p>
            For other information and nicely formatted{" "}
            <span className="discord-role">@mentions</span>, please run{" "}
            <code>/gentei info</code> in your server.
          </p>
        </div>
        <div>
          <h2 className="title is-4">Memberships</h2>
          <p className="mb-1">
            Settings that are hard to configure using slash commands can be
            edited below.
          </p>
          {membershipManagers}
        </div>
        <div className="mt-4">
          <div className="control has-text-centered">
            <Link
              className="button is-primary is-medium"
              to={{
                pathname: "/app/enroll",
                search: new URLSearchParams({ server: guildID }).toString(),
              }}
            >
              Add/Remove Memberships
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}
