import { RiGithubFill } from "react-icons/ri";
import { Link } from "react-router-dom";
import membershipDiagram from "../assets/img/01.png";
import securityKeibiRobot from "../assets/img/security_keibi_robot.png";
import Footer from "../components/Footer";

function Home() {
  return (
    <div className="home">
      <section className="hero is-medium is-primary is-bold">
        <div className="hero-body">
          <div className="container">
            <h1 className="title">Gentei / 限定</h1>
            <h2 className="subtitle">VTuber channel membership verification</h2>
          </div>
        </div>
        <div className="hero-foot has-text-centered pb-5">
          <Link to="/app">
            <button className="button is-link is-size-4">
              Enroll a Community / Validate Membership
            </button>
          </Link>
        </div>
      </section>
      <section className="section"></section>
      <section className="section" role="main">
        <div className="container">
          <div className="columns">
            <div className="column">
              <div className="content has-text-centered">
                <img src={membershipDiagram} alt="scuffed membership diagram" />
              </div>
            </div>
            <div className="column">
              <div className="content">
                <h3>Free membership verification for fans</h3>
                <p>
                  Administrators of fan communities and Discord servers for
                  YouTube channels no longer need to resort to asking for
                  regular screenshots to verify YouTube channel memberships.
                </p>
                <p>
                  Gentei runs a membership verification process on their behalf
                  - this verifies membership via the YouTube API on a regular
                  basis! All for free, no catch.
                </p>
                <p>
                  For technical and administrative documentation, see{" "}
                  <a
                    href="https://docs.gentei.tindabox.net"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    https://docs.gentei.tindabox.net
                  </a>
                  !
                </p>
              </div>
            </div>
          </div>
          <div className="columns">
            <div className="column">
              <div className="content">
                <h3>Discord role assignment bot</h3>
                <p>
                  Take advantage of a Discord bot that can automatically assign
                  and unassign roles to Discord server users. Please use the
                  "Enroll a Community" button above, if interested!
                </p>
                <p>
                  For more info on the bot, please see{" "}
                  <a
                    href="https://docs.gentei.tindabox.net/bot/"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    the subpage in the documentation
                  </a>
                  .
                </p>
              </div>
            </div>
            <div className="column">
              <div className="content has-text-centered">
                <img
                  src={securityKeibiRobot}
                  style={{ maxHeight: 180 }}
                  alt="security robot"
                />
              </div>
            </div>
          </div>
          <div className="columns">
            <div className="column">
              <div className="content has-text-centered">
                <RiGithubFill size={180} />
              </div>
            </div>
            <div className="column">
              <div className="content">
                <h3>Open Source</h3>
                <p>
                  Gentei is an open source, AGPLv3-licensed SaaS project both
                  hosted on GitHub and deployed straight from the project for
                  transparency.
                </p>
                <p>
                  To check out the code and infrastructure, see{" "}
                  <a href="https://github.com/member-gentei">
                    https://github.com/member-gentei
                  </a>
                  .
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>
      <Footer />
    </div>
  );
}

export default Home;
