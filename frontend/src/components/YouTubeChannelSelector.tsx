import { Autocomplete, TextField } from "@mui/material";
import React, { useEffect, useState } from "react";
import { LoadState } from "../lib/lib";
import { useTalents } from "../stores/TalentStore";
import {
  Select,
  Option,
  Grid,
  Input,
  FormHelperText,
  FormControl,
  Button,
} from "@mui/joy";

interface YouTubeChannelSelectorProps {
  selectedChannels: string[];
  addChannel: (channelID: string) => void;
}

const initialHelpText =
  "If you can't find a talent's YouTube channel by name, you can add their channel by URL.";

const reYouTubeChannelID = new RegExp(
  "^https://(www.)?youtube.com/channel/(UC.{22})$",
);

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
      <Grid container sx={{ flexGrow: 1 }}>
        <Grid sm={12} md="auto" lg="auto">
          <Select
            defaultValue="name"
            onChange={(_, value) => {
              setChannelInputType(
                (value || "").toLowerCase() as "name" | "url",
              );
            }}
          >
            <Option value="name">Name</Option>
            <Option value="url">URL</Option>
          </Select>
        </Grid>
        <Grid sm={12} md lg>
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
        </Grid>
      </Grid>
      {channelInputType === "url" ? (
        <Grid container sx={{ mt: 1 }}>
          <Grid sm={0} md={2} lg={1}></Grid>
          <Grid sm md lg sx={{ textAlign: "center" }}>
            <Button type="submit" size="lg" disabled={!inputValid}>
              Add channel
            </Button>
          </Grid>
        </Grid>
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
    <FormControl>
      <Autocomplete
        disablePortal
        options={talentsByName}
        size="small"
        value={selectedChannel}
        onChange={(_, value) => setSelectedChannel(value)}
        renderInput={(params) => <TextField {...params} label="Channel name" />}
      />
      <FormHelperText>{initialHelpText}</FormHelperText>
    </FormControl>
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
  const [neverTouched, setNeverTouched] = useState(true);
  useEffect(() => {
    if (channelURL === "") {
      setValid(undefined);
      setHelpText(initialHelpText);
    }
  }, [channelURL, setValid, setHelpText]);
  return (
    <FormControl>
      <Input
        type="url"
        error={!valid && !neverTouched}
        placeholder="https://www.youtube.com/channel/UC9ruVYPv7yJmV0Rh0NKA-Lw"
        value={channelURL}
        onChange={(e) => {
          setNeverTouched(false);
          const valid = e.target.value.match(reYouTubeChannelID);
          if (!valid) {
            setValid(false);
            setHelpText(
              "Please input a YouTube channel URL. See the placeholder text for an example.",
            );
          } else {
            setValid(true);
            setHelpText("");
          }
          setChannelURL(e.target.value);
        }}
        required
      />
      <FormHelperText>{helpText}</FormHelperText>
    </FormControl>
  );
}
