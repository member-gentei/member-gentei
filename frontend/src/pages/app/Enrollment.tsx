import { Step, StepLabel, Stepper } from "@mui/material";
import { Link, useLocation } from "react-router-dom";
import { GuildContainer } from "../../stores/GuildStore";
import EnrollServer from "./Enrollment/EnrollServer";
import SelectTalents from "./Enrollment/SelectTalents";

enum StepEnum {
  ENROLL_SERVER,
  SELECT_TALENTS,
  CONFIGURE,
}

export default function Enrollment() {
  const location = useLocation();
  const currentStep = evalStep(new URLSearchParams(location.search));
  return (
    <section className="section">
      <div className="container">
        <nav className="breadcrumb">
          <ul>
            <li>
              <Link to="/app">Home</Link>
            </li>
            <li className="is-active">
              <Link to="#">Server Enrollment</Link>
            </li>
          </ul>
        </nav>
        <h1 className="title">Server Enrollment</h1>
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
      </div>
    </section>
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
      <div className="mt-4">
        <Component />
      </div>
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
