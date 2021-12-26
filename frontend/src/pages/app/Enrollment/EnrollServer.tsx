import { Fragment, ReactNode } from "react";
import { SiDiscord } from "react-icons/si";
import { useLocation } from "react-router-dom";
import {
  useDiscordBotURL,
  useDiscordState,
} from "../../../components/LoginURL";
import { useEnroll } from "../../../stores/EnrollStore";

export default function EnrollServer() {
  const botURL = useDiscordBotURL();
  const discordState = useDiscordState();
  const search = new URLSearchParams(useLocation().search);
  const [store, actions] = useEnroll();
  if (!botURL) {
    return (
      <div className="has-text-centered">
        <span className="spinner"></span>
      </div>
    );
  }
  let enrollTop: ReactNode;
  if (search.has("code")) {
    if (!!store.submitError) {
      enrollTop = (
        <div className="columns is-mobile is-centered">
          <div className="column is-three-quarters-tablet is-half-desktop is-half-widescreen is-half-fullhd">
            <div className="message is-danger">
              <div className="message-header">Error adding bot</div>
              <div className="message-body">{store.submitError}</div>
            </div>
          </div>
        </div>
      );
    }
    actions.verifySubmit(search, discordState);
  } else {
    enrollTop = (
      <strong>
        After adding the bot, you will be redirected back to this page to prove
        that you, specifically, can manage that server!
      </strong>
    );
  }
  return (
    <Fragment>
      <h2 className="title is-4">Enroll Server</h2>
      <p className="content">
        Enroll your server with <code>gentei-bouncer#9835</code> to enable
        membership management.
      </p>
      <div className="content has-text-centered">
        <div className="mb-2">{enrollTop}</div>
        <a className="button is-primary spin-hover" href={botURL}>
          <span className="icon-text">Invite gentei-bouncer#9835</span>
          <span className="icon spin-me slow">
            <SiDiscord />
          </span>
        </a>
      </div>
    </Fragment>
  );
}
