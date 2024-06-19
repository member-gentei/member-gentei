import React, { Fragment, ReactNode, useEffect } from "react";
import { SiDiscord } from "react-icons/si";
import {
  Link as RouterLink,
  Navigate,
  useSearchParams,
} from "react-router-dom";
import Footer from "../components/Footer";
import { useDiscordLoginURL, useYouTubeLoginURL } from "../components/LoginURL";
import { LoadState } from "../lib/lib";
import { useUser } from "../stores/UserStore";
import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  CircularProgress,
  Container,
  Grid,
  Stack,
  Typography,
} from "@mui/joy";

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
            <>
              The login code given to us by Discord can't be used - this usually
              means it expired. Please try again!
            </>
          );
          break;
        default:
          explainer = <>Unhandled error logging in - please try again later.</>;
          break;
      }
      cardContent = (
        <Fragment>
          <CardContent>
            <Alert color="warning">{explainer}</Alert>
          </CardContent>
          <CardActions>
            <Button
              component={RouterLink}
              to="/app"
              variant="outlined"
              color="neutral"
            >
              Back to app home
            </Button>
            <Button component="a" href={loginURL!}>
              Try again
            </Button>
          </CardActions>
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
    <>
      <Grid
        container
        direction="column"
        alignItems="center"
        justifyContent="center"
        mt={4}
        mb={4}
      >
        <Grid xs={12} sm={6}>
          <Card>{cardContent}</Card>
        </Grid>
      </Grid>
      <Footer />
    </>
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
            <>
              <Typography>
                Google says that the token is invalid - please try again.
              </Typography>
              <Typography>
                (This is known to happen when you refresh this page.)
              </Typography>
            </>
          );
          break;
        case "YouTube channel belongs to a different user":
          explainer = (
            <>
              <Typography>
                The YouTube channel you tried to log in with belongs to a
                different Discord user of Gentei.
              </Typography>
              <Typography>
                If you want to associate that channel with your currently logged
                in user, sign in as that other user and remove your account.
              </Typography>
            </>
          );
          break;
        default:
          explainer = (
            <>
              <Typography>
                Unhandled error logging in - please try again later.
              </Typography>
            </>
          );
          break;
      }
      cardContent = (
        <>
          <CardContent>
            <Alert color="warning">{explainer}</Alert>
          </CardContent>
          <CardActions>
            <Button component={RouterLink} to="/app">
              Back to app home
            </Button>
            <Button component="a" href={loginURL!}>
              Try again
            </Button>
          </CardActions>
        </>
      );
      break;
    case LoadState.Succeeded:
      return <Navigate to="/app" />;
    default:
      cardContent = (
        <>
          <CardContent>
            <Stack spacing={1}>
              <Typography>Connecting your YouTube account...</Typography>
              <CircularProgress />
            </Stack>
          </CardContent>
        </>
      );
  }
  return (
    <>
      <Grid
        container
        direction="column"
        alignItems="center"
        justifyContent="center"
        mt={4}
        mb={4}
      >
        <Grid xs={12} sm={6}>
          <Card>{cardContent}</Card>
        </Grid>
      </Grid>
      <Footer />
    </>
  );
}
