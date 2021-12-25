import Footer from "../components/Footer";

import membershipDiagram from "../assets/img/01.png";
import securityKeibiRobot from "../assets/img/security_keibi_robot.png";
import monkeyWrench from "../assets/img/monkey_wrench.png";

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
          <a
            href="https://forms.gle/rr4Psqbzz1Nhuqno6"
            target="_blank"
            rel="noopener noreferrer"
          >
            <button className="button is-link is-light is-size-4">
              Enroll a Community
            </button>
          </a>
          <a href="/app">
            <button className="button is-link is-size-4">
              Validate Membership
            </button>
          </a>
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
                  For technical docs, see{" "}
                  <a
                    href="https://docs.member-gentei.tindabox.net"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    https://docs.member-gentei.tindabox.net
                  </a>
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
                    href="https://docs.member-gentei.tindabox.net/Discord/bot"
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
                <img
                  src={monkeyWrench}
                  style={{ maxHeight: 180 }}
                  alt="monkey wrench"
                />
              </div>
            </div>
            <div className="column">
              <div className="content">
                <h3>OpenAPI Integration</h3>
                <p>
                  Fan community administrators can integrate their own tooling
                  and/or existing Discord bots with Gentei's REST API, in lieu
                  of the accompanying Discord bot.
                </p>
                <p>
                  Documentation and authentication details are available in the{" "}
                  <a href="https://docs.member-gentei.tindabox.net/api.html">
                    OpenAPI doc linked here
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
