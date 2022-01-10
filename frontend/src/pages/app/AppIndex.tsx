import React, { Fragment } from "react";
import { RiArrowDropDownLine } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { Link, Outlet } from "react-router-dom";
import logo128 from "../../assets/img/logo-128.png";
import Footer from "../../components/Footer";
import { useDiscordLoginURL } from "../../components/LoginURL";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";

function AppIndex() {
  const actions = useUser()[1];
  actions.getMe();
  return (
    <Fragment>
      <nav className="navbar is-dark">
        <div className="navbar-brand">
          <Link to="/app" className="navbar-item">
            <img src={logo128} alt="Gentei bot logo" />
          </Link>
        </div>
        <div className="navbar-menu">
          <div className="navbar-start"></div>
          <div className="navbar-end">
            <div className="navbar-item">{<AuthButtons />}</div>
          </div>
        </div>
      </nav>
      <Outlet />
      <Footer withYouTubeImage />
    </Fragment>
  );
}

function AuthButtons() {
  const [store, actions] = useUser();
  const loginURL = useDiscordLoginURL();
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
      <div className="dropdown is-right is-hoverable">
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