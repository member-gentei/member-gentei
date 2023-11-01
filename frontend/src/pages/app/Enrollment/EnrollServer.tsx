import { ReactNode } from "react";
import { SiDiscord } from "react-icons/si";
import { useLocation } from "react-router-dom";
import { useDiscordBotURL, useOAuthState } from "../../../components/LoginURL";
import { useEnroll } from "../../../stores/EnrollStore";
import { Box, Button, Stack, Typography } from "@mui/joy";

export default function EnrollServer() {
  const botURL = useDiscordBotURL();
  const discordState = useOAuthState();
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
      <Typography fontWeight="lg">
        After adding the bot, you will be redirected back to this page to prove
        that you, specifically, can manage that server!
      </Typography>
    );
  }
  return (
    <Stack spacing={2}>
      <Typography level="h3">Enroll Server</Typography>
      <Typography>
        Enroll your server by inviting <code>gentei-bouncer#9835</code> to
        enable membership management.
      </Typography>
      <Box sx={{ textAlign: "center" }}>
        <div>{enrollTop}</div>
        <Button
          component="a"
          href={botURL}
          className="spin-hover"
          sx={{ mt: 1, mb: 1 }}
          endDecorator={<SiDiscord className="spin-me" />}
        >
          Invite gentei-bouncer#9835
        </Button>
      </Box>
    </Stack>
  );
}
