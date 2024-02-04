import { Step, StepLabel, Stepper } from "@mui/material";
import { Link as RRDLink, useLocation } from "react-router-dom";
import { GuildContainer } from "../../stores/GuildStore";
import EnrollServer from "./Enrollment/EnrollServer";
import SelectTalents from "./Enrollment/SelectTalents";
import { Box, Breadcrumbs, Link, Stack, Typography } from "@mui/joy";

enum StepEnum {
  ENROLL_SERVER,
  SELECT_TALENTS,
  CONFIGURE,
}

export default function Enrollment() {
  const location = useLocation();
  const currentStep = evalStep(new URLSearchParams(location.search));
  return (
    <Box component="section" sx={{ mb: 2 }}>
      <Stack spacing={2}>
        <Breadcrumbs>
          <RRDLink to="/app">
            <Link>Home</Link>
          </RRDLink>
          <Typography>Server Enrollment</Typography>
        </Breadcrumbs>
        <Typography level="h2">Server Enrollment</Typography>
        <Stepper activeStep={currentStep}>
          <Step>
            <StepLabel>Invite the bot</StepLabel>
          </Step>
          <Step>
            <StepLabel>Pick talent memberships(s)</StepLabel>
          </Step>
          <Step>
            <StepLabel>Configure server + roles</StepLabel>
          </Step>
        </Stepper>
        <hr />
        <EnrollmentStep current={currentStep} />
      </Stack>
    </Box>
  );
}

function EnrollmentStep({ current }: { current: StepEnum }) {
  let Component: () => JSX.Element;
  switch (current) {
    case StepEnum.ENROLL_SERVER:
      Component = EnrollServer;
      break;
    case StepEnum.SELECT_TALENTS:
      Component = SelectTalents;
      break;
    default:
      return <div className="configure-server"></div>;
  }
  return (
    <GuildContainer>
      <Component />
    </GuildContainer>
  );
}

function evalStep(search: URLSearchParams): StepEnum {
  const serverID = search.get("server");
  if (!serverID) {
    return StepEnum.ENROLL_SERVER;
  }
  return StepEnum.SELECT_TALENTS;
}
