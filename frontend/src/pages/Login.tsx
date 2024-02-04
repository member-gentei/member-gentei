import React, { Fragment, ReactNode, useEffect } from "react";
import { SiDiscord } from "react-icons/si";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Footer from "../components/Footer";
import { useDiscordLoginURL, useYouTubeLoginURL } from "../components/LoginURL";
import { LoadState } from "../lib/lib";
import { useUser } from "../stores/UserStore";
import { Typography } from "@mui/joy";

export function LoginDiscord() {
  const [search] = useSearchParams();
  const [store, actions] = useUser();
  const loginURL = useDiscordLoginURL();
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

export function LoginYouTube() {
  const [search] = useSearchParams();
  const [store, actions] = useUser();
  const loginURL = useYouTubeLoginURL();
  useEffect(() => {
    if (store.youtubeLogin !== LoadState.NotStarted) {
      return;
    }
    const code = search.get("code");
    const state = search.get("state");
    if (!!(code && state)) {
      actions.loginYouTube(code, state);
    }
  });
  let cardContent: ReactNode;
  switch (store.youtubeLogin) {
    case LoadState.Failed:
      let explainer: ReactNode;
      switch (store.youtubeLoginError?.error) {
        case "invalid_grant":
          explainer = (
            <Fragment>
              <Typography>
                Google says that the token is invalid - please try again.
              </Typography>
              <Typography>
                (This is known to happen when you refresh this page.)
              </Typography>
            </Fragment>
          );
          break;
        case "YouTube channel belongs to a different user":
          explainer = (
            <Fragment>
              <Typography>
                The YouTube channel you tried to log in with belongs to a
                different Discord user of Gentei.
              </Typography>
              <Typography>
                If you want to associate that channel with your currently logged
                in user, sign in as that other user and remove your account.
              </Typography>
            </Fragment>
          );
          break;
        default:
          explainer = (
            <Fragment>
              <Typography>
                Unhandled error logging in - please try again later.
              </Typography>
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
          <span>Connecting YouTube...</span>
          <div className="mt-2">
            <span className="spin icon m-auto">
              <span className="spinner"></span>
            </span>
          </div>
        </Fragment>
      );
  }
  return (
    <section className="section">
      <div className="container">
        <div className="columns">
          <div className="column">
            <div className="card">
              <div className="card-content has-text-centered">
                {cardContent}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
