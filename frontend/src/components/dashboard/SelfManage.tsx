import classNames from "classnames";
import { ReactNode, useEffect, useState } from "react";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";
import YouTubeLogin from "./YouTubeLogin";

export default function SelfManage() {
  const [store] = useUser();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  useEffect(() => {
    if (store.user === undefined && showDeleteModal) {
      setShowDeleteModal(false);
    }
  }, [store.user, showDeleteModal, setShowDeleteModal]);
  if (store.user === undefined) {
    return null;
  }
  const user = store.user!;
  let youTubeChannelElement: ReactNode;
  if (user.YouTube.ID !== "") {
    const channelURL = `https://youtube.com/channel/${user.YouTube.ID}`;
    youTubeChannelElement = (
      <a href={channelURL} target="_blank" rel="noreferrer">
        {channelURL}
      </a>
    );
  } else {
    youTubeChannelElement = <code>(no channel connected)</code>;
  }
  return (
    <div className="container mt-6">
      <h2 className="title is-3">Self management</h2>
      <p className="my-2">
        This table describes attributes that we use to determine role
        eligibility. Before reaching out for help with role assignment issues,
        please make sure that this is all correct.
      </p>
      <table className="table is-striped">
        <tbody>
          <tr>
            <th>Discord user</th>
            <td>
              <code>{user.FullName}</code>
            </td>
          </tr>
          <tr>
            <th>YouTube channel</th>
            <td>{youTubeChannelElement}</td>
          </tr>
          <tr>
            <th>Actions</th>
            <td>
              <div className="columns">
                <div className="column has-text-centered">
                  <div className="mb-1">
                    Connect a {user.YouTube.ID === "" ? "" : "different "}
                    YouTube account
                  </div>
                  <YouTubeLogin className="" />
                </div>
                <div className="column has-text-centered">
                  <div className="mb-1">
                    Disconnect YouTube account <br />
                    (revokes Discord roles)
                  </div>
                  <button
                    className="button is-danger is-light"
                    disabled={user.YouTube.ID === ""}
                  >
                    Disconnect
                  </button>
                </div>
                <div className="column has-text-centered">
                  <div className="mb-1">
                    Remove your accounts from Gentei and revoke Discord roles
                  </div>
                  <button
                    className="button is-danger"
                    onClick={() => setShowDeleteModal(true)}
                  >
                    Remove
                  </button>
                </div>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <DeleteModal
        show={showDeleteModal}
        hide={() => setShowDeleteModal(false)}
      />
    </div>
  );
}

interface DeleteModalProps {
  show: boolean;
  hide: () => void;
}

function DeleteModal({ show, hide }: DeleteModalProps) {
  const [store, actions] = useUser();
  const className = classNames("modal", { "is-active": show });
  // press escape to hide the modal
  useEffect(() => {
    const keyDownListener = (e: KeyboardEvent) => {
      if (e.code === "Escape") {
        hide();
      }
    };
    document.addEventListener("keydown", keyDownListener);
    return () => {
      document.removeEventListener("keydown", keyDownListener);
    };
  });
  const deleteButtonClassNames = classNames("button", "is-danger", {
    "is-loading": store.remove === LoadState.Started,
  });
  return (
    <div className={className}>
      <div className="modal-background" onClick={hide}></div>
      <div className="modal-card">
        <header className="modal-card-head">
          <p className="modal-card-title">Account deletion</p>
          <button className="delete" onClick={hide}></button>
        </header>
        <section className="modal-card-body">
          <div className="content">
            <p>
              Sorry to see you go! Here's exactly what will happen if you
              confirm account deletion:
            </p>
            <ol>
              <li>
                You will be removed from all Discord roles that{" "}
                <code>gentei-bouncer#9835</code> has assigned. Participating
                Discord servers may or may not have different members-only roles
                that you can use.
              </li>
              <li>
                We will revoke all access tokens and delete all information
                about you <em>except</em> your Discord user ID. It will still be
                present in audit logs, per-server access control lists for
                running management commands, and ~weekly discarded database
                backups.
              </li>
              <li>
                If you added Gentei to a Discord server,{" "}
                <strong>the bot will stay in that server until kicked</strong>.
                Please kick the bot to remove it from your server.
              </li>
            </ol>
          </div>
        </section>
        <footer className="modal-card-foot">
          <button
            className={deleteButtonClassNames}
            onClick={() => actions.remove()}
          >
            Delete my account
          </button>
          <button className="button">Cancel</button>
        </footer>
      </div>
    </div>
  );
}
