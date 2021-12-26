import { useState } from "react";
import { DISCORD_CLIENT_ID, REDIRECT_URI } from "../lib/lib";

export function useDiscordLoginURL(): string | void {
  // discordState should be a pretty random string
  const [discordState, setDiscordState] = useState(
    localStorage.getItem("state")
  );
  if (!discordState) {
    let array = new Uint8Array(16);
    crypto.getRandomValues(array);
    let hexState = "";
    array.forEach((b) => (hexState += b.toString(16).padStart(2, "0")));
    localStorage.setItem("state", hexState);
    setDiscordState(hexState);
    return;
  }
  const loginParams = new URLSearchParams({
    client_id: DISCORD_CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    response_type: "code",
    scope: "identify guilds",
    state: discordState,
  });
  return `https://discord.com/api/oauth2/authorize?${loginParams.toString()}`;
}
