import React, { Fragment, ReactNode, useEffect } from "react";
import { SiDiscord } from "react-icons/si";
import { Link, Navigate, useLocation } from "react-router-dom";
import Footer from "../components/Footer";
import { useDiscordLoginURL } from "../components/LoginURL";
import { LoadState } from "../lib/lib";
import { useUser } from "../stores/UserStore";

export function LoginDiscord() {
  const location = useLocation();
  const [store, actions] = useUser();
  const loginURL = useDiscordLoginURL();
  const search = new URLSearchParams(location.search);
  useEffect(() => {
    if (store.discordLogin !== LoadState.NotStarted) {
      return;
    }
    const code = search.get("code");
    const state = search.get("state");
    if (!!(code && state)) {
      actions.loginDiscord(code, state);
    }
  });
  let cardContent: ReactNode;
  let columnsClasses = "";
  let columnClasses =
    "is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-quarter-fullhd";
  switch (store.discordLogin) {
    case LoadState.Failed:
      let explainer: ReactNode;
      columnClasses = "";
      columnsClasses = "mx-6";
      switch (store.discordLoginError?.error_description) {
        case `Invalid "code" in request.`:
          explainer = (
            <Fragment>
              The login code given to us by Discord can't be used - this usually
              means it expired. Please try again!
            </Fragment>
          );
          break;
        default:
          explainer = (
            <Fragment>
              Unhandled error logging in - please try again later.
            </Fragment>
          );
          break;
      }
      cardContent = (
        <Fragment>
          <div className="notification is-danger">{explainer}</div>
          <div className="buttons is-centered">
            <Link to="/app" className="button is-secondary">
              Back to app home
            </Link>
            <a className="button is-primary" href={loginURL!}>
              Try again
            </a>
          </div>
        </Fragment>
      );
      break;
    case LoadState.Succeeded:
      return <Navigate to="/app" />;
    default:
      cardContent = (
        <Fragment>
          <span>Logging you in with Discord...</span>
          <div className="mt-2">
            <span className="spin icon m-auto">
              <SiDiscord />
            </span>
          </div>
        </Fragment>
      );
  }
  return (
    <section className="section">
      <div className="container">
        <div className={`columns is-mobile is-centered mt-4 ${columnsClasses}`}>
          <div className={`column ${columnClasses}`}>
            <div className="card">
              <div className="card-content has-text-centered">
                {cardContent}
              </div>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </section>
  );
}
