import { IoPersonOutline } from "react-icons/io5";
import { LoadState } from "../lib/lib";
import { useTalents } from "../stores/TalentStore";
import {
  AspectRatio,
  Button,
  Card,
  CardActions,
  CardContent,
  Divider,
  Stack,
  Typography,
} from "@mui/joy";

interface TalentCardProps {
  channelID: string;
  error?: boolean;
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
  let cardActions = null;
  if (onDelete !== undefined) {
    cardActions = (
      <CardActions>
        <Button color="danger" onClick={onDelete}>
          Remove
        </Button>
      </CardActions>
    );
  }
  const talent = store.talentsByID[channelID];
  if (!talent) {
    return (
      <Card sx={{ textAlign: "center" }}>
        <AspectRatio ratio="1/1">
          <IoPersonOutline size={128} />
        </AspectRatio>
        <CardContent>
          <Stack spacing={1}>
            <em>
              <a href={channelURL}>
                new channel <br />{" "}
                <span className="is-size-7" style={{ whiteSpace: "nowrap" }}>
                  ({channelID})
                </span>
              </a>
            </em>
            <Divider />
            <Typography>
              New channels are processed after submission.
            </Typography>
          </Stack>
        </CardContent>
        <div className="card-content is-clipped">
          <div className="content has-text-centered"></div>
          <div className="content"></div>
        </div>
        {cardActions}
      </Card>
    );
  }
  let talentThumbnailNode;
  if (talent.Thumbnail !== "") {
    talentThumbnailNode = (
      <img src={talent.Thumbnail} alt="channel thumbnail" />
    );
  } else {
    talentThumbnailNode = <IoPersonOutline size={128} />;
  }
  return (
    <Card sx={{ textAlign: "center" }}>
      <AspectRatio ratio="1/1">{talentThumbnailNode}</AspectRatio>
      <CardContent>
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
      </CardContent>
      {cardActions}
    </Card>
  );
}
