import classNames from "classnames";
import React, { Fragment, useState } from "react";
import { RiArrowDropDownLine } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { Link, Outlet } from "react-router-dom";
import logo128 from "../../assets/img/logo-128.png";
import Footer from "../../components/Footer";
import { useDiscordLoginURL } from "../../components/LoginURL";
import { LoadState, useWindowSize } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";

const discordInviteURL = "https://discord.gg/xJd9Der";

function AppIndex() {
  const actions = useUser()[1];
  const [navActive, setNavActive] = useState(false);
  actions.getMe();
  return (
    <Fragment>
      <nav className="navbar is-dark">
        <div className="navbar-brand">
          <Link to="/app" className="navbar-item">
            <img src={logo128} alt="Gentei bot logo" />
          </Link>
          <a
            href="#"
            role="button"
            className={classNames("navbar-burger", { "is-active": navActive })}
            onClick={() => setNavActive((v) => !v)}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>
        <div className={classNames("navbar-menu", { "is-active": navActive })}>
          <div className="navbar-start"></div>
          <div className="navbar-end">
            <div className="navbar-item">{<AuthButtons />}</div>
          </div>
        </div>
      </nav>
      <section className="section pt-4 pb-0">
        <div
          className="message is-info"
          style={{
            maxWidth: 640,
            margin: "auto",
          }}
        >
          <div className="message-header">
            <p>Soft launch preview</p>
          </div>
          <div className="message-body">
            <p>
              Missing features and documentation will be slowly turned on over
              the next couple of weeks. Please see the{" "}
              <code>#gentei-announce</code> channel in the{" "}
              <a href={discordInviteURL}>Hololive Creators Club Discord</a>{" "}
              server for updates and instructions - both for users and community
              admins - while this message visible.
            </p>
          </div>
        </div>
      </section>
      <Outlet />
      <Footer withYouTubeImage />
    </Fragment>
  );
}

function AuthButtons() {
  const [store, actions] = useUser();
  const loginURL = useDiscordLoginURL();
  const windowSize = useWindowSize();
  if (store.userLoad <= LoadState.Started || !loginURL) {
    return (
      <progress className="progress is-small" max="100">
        69%
      </progress>
    );
  }
  const logout: React.MouseEventHandler<HTMLAnchorElement> = (e) => {
    e.preventDefault();
    actions.logout();
  };
  if (store.userLoad === LoadState.Succeeded && !!store.user) {
    const user = store.user!;
    const avatarURL = `https://cdn.discordapp.com/avatars/${user.ID}/${user.AvatarHash}.webp?size=64`;
    return (
      <div
        className={classNames("dropdown is-hoverable", {
          "is-right": windowSize.width >= 1024,
        })}
      >
        <div className="dropdown-trigger">
          <button className="button is-black outlined">
            <span>
              <figure className="image avatar is-square">
                <img
                  src={avatarURL}
                  alt={`Discord avatar for ${user.FullName}`}
                  className="is-rounded"
                />
              </figure>
            </span>
            <span className="icon">
              <RiArrowDropDownLine size="2em" />
            </span>
          </button>
        </div>
        <div className="dropdown-menu">
          <div className="dropdown-content">
            <span className="dropdown-item">{user.FullName}</span>
            <hr className="dropdown-divider" />
            <a className="dropdown-item" href="/logout" onClick={logout}>
              Sign out
            </a>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="buttons">
      <a className="button is-primary" href={loginURL}>
        <span className="icon-text">
          <span>Register / Sign in with Discord</span>
          <span className="icon">
            <SiDiscord />
          </span>
        </span>
      </a>
    </div>
  );
}

export default AppIndex;
