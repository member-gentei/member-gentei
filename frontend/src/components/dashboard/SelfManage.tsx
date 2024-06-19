import classNames from "classnames";
import { ReactNode, useEffect, useState } from "react";
import { LoadState } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";
import YouTubeLogin from "./YouTubeLogin";
import {
  Box,
  Button,
  DialogContent,
  DialogTitle,
  Grid,
  Link,
  Modal,
  ModalClose,
  ModalDialog,
  Sheet,
  Stack,
  Table,
  Typography,
} from "@mui/joy";

export default function SelfManage() {
  const [store, actions] = useUser();
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
      <Link href={channelURL} target="_blank" rel="noreferrer">
        {channelURL}
      </Link>
    );
  } else {
    youTubeChannelElement = <code>(no channel connected)</code>;
  }
  return (
    <Stack spacing={1}>
      <Typography level="h2">Self management</Typography>
      <Typography>
        This table describes attributes that we use to determine role
        eligibility. Before reaching out for help with role assignment issues,
        please make sure that this is all correct.
      </Typography>
      <Sheet sx={{ maxWidth: 1100 }}>
        <Table
          borderAxis="bothBetween"
          sx={{
            "& th": { width: "10rem" },
            overflow: "hidden",
          }}
        >
          <tbody>
            <tr>
              <th scope="row">Discord user</th>
              <td>
                <code>{user.FullName}</code>
              </td>
            </tr>
            <tr>
              <th scope="row">YouTube channel</th>
              <td>{youTubeChannelElement}</td>
            </tr>
            <tr>
              <th scope="row">Actions</th>
              <td>
                <Grid container columnSpacing={3} sx={{ textAlign: "center" }}>
                  <Grid sm={6} md={4}>
                    <div className="mb-1">
                      Connect a {user.YouTube.ID === "" ? "" : "different "}
                      YouTube account
                    </div>
                    <YouTubeLogin className="" />
                  </Grid>
                  <Grid sm={6} md={4}>
                    <div className="mb-1">
                      Disconnect YouTube account <br />
                      (revokes Discord roles)
                    </div>
                    <Button
                      color="danger"
                      variant="soft"
                      disabled={user.YouTube.ID === ""}
                      loading={store.disconnect === LoadState.Started}
                      onClick={() => actions.disconnectYouTube()}
                    >
                      Disconnect
                    </Button>
                  </Grid>
                  <Grid md={4}>
                    <div className="mb-1">
                      Remove your accounts from Gentei and revoke Discord roles
                    </div>
                    <Button
                      color="danger"
                      onClick={() => setShowDeleteModal(true)}
                    >
                      Remove
                    </Button>
                  </Grid>
                </Grid>
              </td>
            </tr>
          </tbody>
        </Table>
      </Sheet>
      <DeleteModal
        show={showDeleteModal}
        hide={() => setShowDeleteModal(false)}
      />
    </Stack>
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
    <Modal
      className={className}
      open={show}
      onClose={hide}
      sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
    >
      <ModalDialog color="danger">
        <ModalClose />
        <DialogTitle>Account deletion</DialogTitle>
        <DialogContent>
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
                  Discord servers may or may not have different members-only
                  roles that you can use.
                </li>
                <li>
                  We will revoke all access tokens and delete all information
                  about you <em>except</em> your Discord user ID. It will still
                  be present in audit logs, per-server access control lists for
                  running management commands, and ~weekly discarded database
                  backups.
                </li>
                <li>
                  If you added Gentei to a Discord server,{" "}
                  <strong>the bot will stay in that server until kicked</strong>
                  . Please kick the bot to remove it from your server.
                </li>
              </ol>
            </div>
          </section>
          <Stack direction="row" spacing={1}>
            <Button color="danger" onClick={() => actions.remove()}>
              Delete my account
            </Button>
            <Button color="neutral" onClick={hide}>
              Cancel
            </Button>
          </Stack>
        </DialogContent>
      </ModalDialog>
    </Modal>
  );
}
