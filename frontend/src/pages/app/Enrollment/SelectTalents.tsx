import classNames from "classnames";
import React, { Fragment, useEffect, useState } from "react";
import { SiDiscord } from "react-icons/si";
import { useNavigate, useSearchParams } from "react-router-dom";
import TalentCard from "../../../components/TalentCard";
import YouTubeChannelSelector from "../../../components/YouTubeChannelSelector";
import { LoadState } from "../../../lib/lib";
import { useGuild } from "../../../stores/GuildStore";

export default function SelectTalents() {
  const [search, setSearch] = useSearchParams();
  const [store, actions] = useGuild();
  const [ready, setReady] = useState(false);
  actions.load(search.get("server")!);
  // run once after guildState changes
  useEffect(() => {
    if (store.guildState <= LoadState.Started) {
      return;
    }
    const talentIDs = store.guild!.TalentIDs || [];
    if (search.getAll(talentGetParam).length === 0 && talentIDs.length > 0) {
      talentIDs.forEach((talentID) => {
        search.append(talentGetParam, talentID);
      });
      setSearch(search, { replace: true });
    }
    setReady(true);
    // eslint-disable-next-line
  }, [setReady, store.guildState]);
  if (store.guildState <= LoadState.Started || !ready) {
    return (
      <div className="has-text-centered">
        <SiDiscord className="spin mt-4" size={24} />
        <div>Loading Discord server info...</div>
      </div>
    );
  }
  if (store.guildState === LoadState.Failed) {
    return (
      <div className="columns is-mobile is-centered">
        <div className="column is-three-quarters-tablet is-half-desktop is-half-widescreen is-half-fullhd">
          <div className="message is-danger">
            <div className="message-header">Error adding bot</div>
            <div className="message-body">{store.guildError}</div>
          </div>
        </div>
      </div>
    );
  }
  return <SelectTalentsInner />;
}

const talentGetParam = "talent";

function SelectTalentsInner() {
  const navigate = useNavigate();
  const [search, setSearch] = useSearchParams();
  const [store, actions] = useGuild();
  const [selectedTalentIDs, setSelectedTalentIDs] = useState<string[]>(
    search.getAll(talentGetParam)
  );
  useEffect(() => {
    if (store.saveTalentsState === LoadState.Succeeded) {
      navigate(`/app/server/${store.guild!.ID}`);
    }
  }, [store.guild, store.saveTalentsState, navigate]);
  useEffect(() => {
    const talentParams = search.getAll(talentGetParam);
    if (talentParams.length !== selectedTalentIDs.length) {
      setSelectedTalentIDs(talentParams);
    }
  }, [search, selectedTalentIDs]);
  const talentCards = selectedTalentIDs.map((channelID) => {
    return (
      <TalentCard
        key={`tc-${channelID}`}
        channelID={channelID}
        error={!!store.guildError?.talents?.[channelID]}
        onDelete={() => {
          // recreate param
          search.delete(talentGetParam);
          selectedTalentIDs.forEach((v) => {
            if (v !== channelID) {
              search.append(talentGetParam, v);
            }
          });
          setSearch(search);
        }}
      />
    );
  });
  let saveDisabled =
    selectedTalentIDs.length === 0 ||
    store.saveTalentsState === LoadState.Failed;
  let errorNode = null;
  if (selectedTalentIDs.length === 0) {
    errorNode = (
      <p className="help">
        Servers must be configured to track at least one membership.
      </p>
    );
  }
  if (store.saveTalentsState === LoadState.Failed) {
    console.log(store.saveTalentsError);
    if (store.saveTalentsError?.talents !== undefined) {
      if (
        errorTalentsRemoved(store.saveTalentsError.talents, selectedTalentIDs)
      ) {
        saveDisabled = false;
      } else {
        const lis = Object.entries(store.saveTalentsError.talents).map(
          ([talentID, msg]) => (
            <li key={`${talentID}-error`}>
              <span className="has-text-weight-bold">{talentID}</span>: {msg}
            </li>
          )
        );
        errorNode = (
          <div className="message is-danger">
            <div className="message-body">
              Error(s) adding talents. Remove them above before proceeding.
              <ul>{lis}</ul>
            </div>
          </div>
        );
      }
    } else {
      errorNode = <p className="help">{store.saveTalentsError?.message}</p>;
    }
  }
  return (
    <Fragment>
      <h2 className="title is-4">Select Talent(s)</h2>
      <p className="content">
        Please select or add YouTube channels whose memberships should be
        monitored for the <strong>{store.guild?.Name}</strong> server.
      </p>
      <div className="content">
        <YouTubeChannelSelector
          selectedChannels={selectedTalentIDs}
          addChannel={(channelID) => {
            if (selectedTalentIDs.indexOf(channelID) === -1) {
              search.append(talentGetParam, channelID);
              setSearch(search, { replace: true });
            }
          }}
        />
      </div>
      <div className="content is-flex is-flex-wrap-wrap is-justify-content-center is-align-items-center">
        {talentCards}
      </div>
      <div className="content">
        <div className="columns is-centered">
          <div className="column is-half">{errorNode}</div>
        </div>
        <div className="has-text-centered">
          <button
            className={`button is-primary is-medium ${classNames({
              "is-loading": store.saveTalentsState === LoadState.Started,
            })}`}
            disabled={saveDisabled}
            onClick={(e) => {
              e.preventDefault();
              actions.saveTalentChannels(store.guild!.ID, selectedTalentIDs);
            }}
          >
            Save YouTube Channels
          </button>
        </div>
      </div>
    </Fragment>
  );
}

function errorTalentsRemoved(
  errors: {
    [key: string]: string | undefined;
  },
  selectedTalentIDs: string[]
): boolean {
  for (const talentID of selectedTalentIDs) {
    if (errors[talentID] !== undefined) {
      return false;
    }
  }
  return true;
}
