import { Fragment } from "react";
import AdminServers from "../../components/dashboard/AdminServers";
import DiscordServers from "../../components/dashboard/DiscordServers";
import SelfManage from "../../components/dashboard/SelfManage";
import { YouTubeLoginOverlay } from "../../components/dashboard/YouTubeLogin";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";
import { Box, Button, Grid, Stack } from "@mui/joy";
import { useDiscordLoginURL } from "../../components/LoginURL";
import { SiDiscord } from "react-icons/si";

export default function UserDashboard() {
  const [store] = useUser();
  let serversOrLoginNode;
  if (store.user === undefined || !store.user.YouTube.Valid) {
    serversOrLoginNode = (
      <Stack spacing={4}>
        <Box component="section" sx={{ position: "relative" }}>
          <DiscordServers />
          <YouTubeLoginOverlay />
        </Box>
        {!store.user ? <RegisterSignIn /> : null}
      </Stack>
    );
  } else {
    serversOrLoginNode = (
      <section>
        <DiscordServers />
      </section>
    );
  }
  return (
    <Grid container rowSpacing={2}>
      <Grid xs={12}>
        {serversOrLoginNode}
      </Grid>
      <Grid xs={12}>
        <SelfManage />
      </Grid>
      <Grid xs={12}>
        {!!store.user ? <AdminServers /> : null}
      </Grid>
    </Grid>
  );
}

function RegisterSignIn() {
  const loginURL = useDiscordLoginURL();
  return (
    <Box sx={{ textAlign: "center" }}>
      <Button
        component="a"
        href={loginURL || "#"}
        startDecorator={<SiDiscord className="spin-me" />}
      >
        Register / Sign in with Discord
      </Button>
    </Box>
  )
}