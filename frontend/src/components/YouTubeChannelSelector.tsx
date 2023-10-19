import { Autocomplete, TextField } from "@mui/material";
import classNames from "classnames";
import React, { useEffect, useState } from "react";
import { LoadState } from "../lib/lib";
import { useTalents } from "../stores/TalentStore";

interface YouTubeChannelSelectorProps {
  selectedChannels: string[];
  addChannel: (channelID: string) => void;
}

const initialHelpText =
  "If you can't find a talent's YouTube channel by name, you can add their channel by URL.";

const reYouTubeChannelID = "https://(www.)?youtube.com/channel/(.{24})";

export default function YouTubeChannelSelector({
  addChannel,
}: YouTubeChannelSelectorProps) {
  const [channelInputType, setChannelInputType] = useState<"name" | "url">(
    "name",
  );
  const [channelURL, setChannelURL] = useState("");
  const [inputValid, setInputValid] = useState<boolean>();
  const [state, actions] = useTalents();
  actions.loadAll();
  useEffect(() => {
    setInputValid(undefined);
  }, [channelInputType]);
  if (state.loadAllState <= LoadState.Started) {
    return (
      <div
        className="has-text-centered box m-auto"
        style={{
          maxWidth: 400,
        }}
      >
        <div className="mb-2">
          <span className="spinner mx-auto"></span>
        </div>
        <div>Loading existing channel data...</div>
      </div>
    );
  }
  function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    switch (channelInputType) {
      case "url":
        addChannel(channelURL.match(reYouTubeChannelID)![2]);
        setChannelURL("");
        setInputValid(undefined);
        break;
    }
  }
  return (
    <form onSubmit={onSubmit}>
      <div className="field is-horizontal">
        <div className="field-label is-normal">
          <label className="label">Channel</label>
        </div>
        <div className="field-body">
          <div className="field has-addons">
            <div className="control">
              <span
                className={classNames("select", { "is-success": !!inputValid })}
              >
                <select
                  name="channel-type"
                  value={channelInputType}
                  onChange={(e) => {
                    setChannelInputType(
                      e.target.value.toLowerCase() as "name" | "url",
                    );
                  }}
                >
                  <option value="name">Name</option>
                  <option value="url">URL</option>
                </select>
              </span>
            </div>
            {channelInputType === "name" ? (
              <ChannelNameSelector addChannel={addChannel} />
            ) : (
              <ChannelURLEntryControls
                valid={inputValid}
                setValid={setInputValid}
                channelURL={channelURL}
                setChannelURL={setChannelURL}
              />
            )}
          </div>
        </div>
      </div>
      {channelInputType === "url" ? (
        <div className="field">
          <div className="control has-text-centered">
            <input
              className="button is-primary"
              type="submit"
              value="Add channel"
              disabled={!inputValid}
            />
          </div>
        </div>
      ) : null}
    </form>
  );
}

interface SelectedChannel {
  id: string;
  label: string;
}

interface ChannelNameSelectorProps {
  addChannel: YouTubeChannelSelectorProps["addChannel"];
}

function ChannelNameSelector({ addChannel }: ChannelNameSelectorProps) {
  const state = useTalents()[0];
  const [selectedChannel, setSelectedChannel] =
    useState<SelectedChannel | null>(null);
  useEffect(() => {
    if (selectedChannel !== null) {
      addChannel(selectedChannel.id);
      setSelectedChannel(null);
    }
  }, [addChannel, selectedChannel]);
  const talentsByName: SelectedChannel[] = Object.keys(state.talentsByID).map(
    (channelID) => ({
      id: channelID,
      label: state.talentsByID[channelID]!.Name,
    }),
  );
  return (
    <div className="control is-expanded">
      <Autocomplete
        disablePortal
        options={talentsByName}
        size="small"
        value={selectedChannel}
        onChange={(_, value) => setSelectedChannel(value)}
        renderInput={(params) => <TextField {...params} label="Channel name" />}
      />
      <p className="help">{initialHelpText}</p>
    </div>
  );
}

interface ChannelURLEntryProps {
  valid?: boolean;
  setValid: (valid?: boolean) => void;
  channelURL: string;
  setChannelURL: (channelID: string) => void;
}

function ChannelURLEntryControls({
  valid,
  setValid,
  channelURL,
  setChannelURL,
}: ChannelURLEntryProps) {
  const [helpText, setHelpText] = useState(initialHelpText);
  useEffect(() => {
    if (channelURL === "") {
      setValid(undefined);
      setHelpText(initialHelpText);
    }
  }, [channelURL, setValid, setHelpText]);
  return (
    <div className="control is-expanded">
      <input
        type="url"
        className={classNames("input", {
          "is-danger": valid === undefined ? false : !valid,
          "is-success": valid,
        })}
        placeholder="https://www.youtube.com/channel/UC9ruVYPv7yJmV0Rh0NKA-Lw"
        value={channelURL}
        pattern={reYouTubeChannelID}
        onChange={(e) => {
          if (!e.target.validity.valid) {
            setValid(false);
            setHelpText("Please input a YouTube channel URL.");
          } else {
            setValid(true);
            setHelpText("");
          }
          setChannelURL(e.target.value);
        }}
        required
      />
      <p
        className={classNames("help", {
          "is-danger": valid === undefined ? false : !valid,
          "is-success": valid,
        })}
      >
        {helpText}
      </p>
    </div>
  );
}
