import { IoPersonOutline } from "react-icons/io5";
import { LoadState } from "../lib/lib";
import { useTalents } from "../stores/TalentStore";
import styles from "./TalentCard.module.css";

interface TalentCardProps {
  channelID: string;
  onDelete?: () => void;
  cardClassNames?: string[];
}

export default function TalentCard({
  channelID,
  cardClassNames,
  onDelete,
}: TalentCardProps) {
  const [store, actions] = useTalents();
  actions.loadAll();
  if (store.loadAllState <= LoadState.Started) {
    // TODO: spinners everywhere
  }
  const channelURL = `https://www.youtube.com/channel/${channelID}`;
  let footerNode = null;
  if (onDelete !== undefined) {
    footerNode = (
      <footer className="card-footer">
        <button
          className="card-footer-item button is-danger is-light"
          onClick={onDelete}
        >
          Remove
        </button>
      </footer>
    );
  }
  const talent = store.talentsByID[channelID];
  const cardClassName = `card m-1 ${(cardClassNames || []).join(" ")} ${
    styles.talentCard
  }`;
  if (!talent) {
    return (
      <div className={cardClassName}>
        <div className="card-image">
          <figure className="image is-128x128">
            <IoPersonOutline size={128} />
          </figure>
        </div>
        <div className="card-content is-clipped">
          <div className="content has-text-centered">
            <em>
              <a href={channelURL}>
                new channel <br />{" "}
                <span className="is-size-7" style={{ whiteSpace: "nowrap" }}>
                  ({channelID})
                </span>
              </a>
            </em>
          </div>
          <div className="content">
            <hr />
            New channels are processed after submission.
          </div>
        </div>
        {footerNode}
      </div>
    );
  }
  return (
    <div className={cardClassName}>
      <div className="card-image">
        <figure className="image is-128x128">
          <img
            className="is-rounded"
            src={talent.Thumbnail}
            alt="channel thumbnail"
          />
        </figure>
      </div>
      <div className="card-content is-clipped">
        <div className="content has-text-centered">
          <strong>
            <a
              href={channelURL}
              target="_blank"
              rel="noreferrer"
              title="Open channel in a new tab"
            >
              {talent.Name}
            </a>
          </strong>
        </div>
      </div>
      {footerNode}
    </div>
  );
}
