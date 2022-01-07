import React from "react";
import { SiDiscord } from "react-icons/si";
import { Link } from "react-router-dom";
import { LoadState } from "../../lib/lib";
import { GuildContainer, useGuild } from "../../stores/GuildStore";
import { useUser } from "../../stores/UserStore";

export default function AdminServers() {
  const [store, actions] = useUser();
  actions.getMe();
  if (store.userLoad <= LoadState.Started) {
    return (
      <div className="section">
        <div className="container">
          <h2 className="title is-3">Administration</h2>
          <div className="has-text-centered`">
            <span className="spinner mx-auto"></span>
          </div>
        </div>
      </div>
    );
  }
  let serverNode;
  if (!!store.user?.ServerAdmin) {
    serverNode = (
      <div className="content is-flex is-flex-wrap-wrap is-justify-content-center">
        {store.user.ServerAdmin.map((guildID) => (
          <AdminServerSnippet key={`${guildID}-admin`} guildID={guildID} />
        ))}
      </div>
    );
  } else {
    serverNode = (
      <p className="content">
        You do not have permission to configure or audit Gentei for any Discord
        server integrations.
      </p>
    );
  }
  return (
    <section className="section">
      <div className="container">
        <h2 className="title is-3">Administration</h2>
        {serverNode}
        <div className="has-text-centered mt-4">
          <Link to="enroll" className="button is-primary spin-hover">
            <span className="icon-text">
              <span>Enroll a new Discord server</span>
              <span className="icon spin-me">
                <SiDiscord />
              </span>
            </span>
          </Link>
        </div>
      </div>
    </section>
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
    iconNode = <img src={guild.Icon} alt="Discord server icon" />;
  } else {
    iconNode = <SiDiscord size={128} />;
  }
  return (
    <div className="box">
      <div className="media">
        <figure className="media-left">
          <p className="image is-128x128">{iconNode}</p>
        </figure>
        <div className="media-content">
          <table className="table">
            <tbody>
              <tr>
                <th>Name</th>
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
                <th>Memberships</th>
                <td>{guild.TalentIDs?.length || 0}</td>
              </tr>
            </tbody>
          </table>
          <Link className="button is-link is-small" to={`server/${guild.ID}`}>
            View/Edit
          </Link>
        </div>
      </div>
    </div>
  );
}
