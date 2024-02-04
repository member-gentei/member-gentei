import { RiGithubFill } from "react-icons/ri";
import { Link as RRDLink } from "react-router-dom";
import membershipDiagram from "../assets/img/01.png";
import securityKeibiRobot from "../assets/img/security_keibi_robot.png";
import Footer from "../components/Footer";
import {
  Box,
  Button,
  Container,
  Grid,
  Link,
  Sheet,
  Stack,
  Typography,
} from "@mui/joy";

function Home() {
  return (
    <div className="home">
      <Sheet
        variant="solid"
        sx={{
          backgroundColor: "rgb(39, 128, 227)",
          padding: 4,
          mb: 12,
        }}
      >
        <Box sx={{ mt: 8, mb: 6 }}>
          <Typography level="h1" sx={{ color: "white" }}>
            Gentei / 限定
          </Typography>
          <Typography level="title-lg" sx={{ color: "white" }}>
            VTuber channel membership verification
          </Typography>
        </Box>
        <Box sx={{ textAlign: "center" }}>
          <RRDLink to="/app">
            <Button size="lg" variant="soft">
              Enroll a Community / Validate Membership
            </Button>
          </RRDLink>
        </Box>
      </Sheet>
      <Container component="section" role="main">
        <Grid container spacing={4} sx={{ mb: 2 }}>
          <Grid xs={6}>
            <Box>
              <img
                src={membershipDiagram}
                alt="scuffed membership diagram"
                style={{ maxWidth: "100%" }}
              />
            </Box>
          </Grid>
          <Grid xs={6}>
            <Stack spacing={2}>
              <Typography level="h3">
                Free membership verification for fans
              </Typography>
              <Typography>
                Administrators of fan communities and Discord servers for
                YouTube channels no longer need to resort to asking for regular
                screenshots to verify YouTube channel memberships.
              </Typography>
              <Typography>
                Gentei runs a membership verification process on their behalf -
                this verifies membership via the YouTube API on a regular basis!
                All for free, no catch.
              </Typography>
              <Typography>
                For technical and administrative documentation, see{" "}
                <Link
                  href="https://docs.gentei.tindabox.net"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  https://docs.gentei.tindabox.net
                </Link>
                !
              </Typography>
            </Stack>
          </Grid>
          <Grid xs={6}>
            <Stack spacing={2}>
              <Typography level="h3">Discord role assignment bot</Typography>
              <Typography>
                Take advantage of a Discord bot that can automatically assign
                and unassign roles to Discord server users. Please use the
                "Enroll a Community" button above, if interested!
              </Typography>
              <Typography>
                For more info on the bot, please see{" "}
                <Link
                  href="https://docs.gentei.tindabox.net/bot/"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  the subpage in the documentation
                </Link>
                !
              </Typography>
            </Stack>
          </Grid>
          <Grid xs={6}>
            <Box textAlign="center">
              <img
                src={securityKeibiRobot}
                style={{ maxHeight: 180 }}
                alt="security robot"
              />
            </Box>
          </Grid>
          <Grid xs={6}>
            <RiGithubFill size={180} />
          </Grid>
          <Grid xs={6}>
            <Stack spacing={2}>
              <Typography level="h3">Open Source</Typography>
              <Typography>
                Gentei is an open source, AGPLv3-licensed SaaS project both
                hosted on GitHub and deployed straight from the project for
                transparency.
              </Typography>
              <Typography>
                To check out the code and infrastructure, see{" "}
                <Link href="https://github.com/member-gentei">
                  https://github.com/member-gentei
                </Link>
                .
              </Typography>
            </Stack>
          </Grid>
        </Grid>
      </Container>
      <Footer />
    </div>
  );
}

export default Home;
