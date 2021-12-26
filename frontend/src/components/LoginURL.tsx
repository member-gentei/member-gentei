import { useState } from "react";
import {
  BOT_REDIRECT_URI,
  DISCORD_BOT_PERMISSIONS,
  DISCORD_CLIENT_ID,
  REDIRECT_URI,
} from "../lib/lib";

export function useDiscordLoginURL(
  additionalScopes: string[] = []
): string | void {
  const discordState = useDiscordState();
  if (!discordState) {
    return;
  }
  const scopes = ["identify", "guilds"].concat(additionalScopes);
  const loginParams = new URLSearchParams({
    client_id: DISCORD_CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    response_type: "code",
    scope: scopes.join(" "),
    state: discordState,
  });
  return `https://discord.com/api/oauth2/authorize?${loginParams.toString()}`;
}

export function useDiscordBotURL(): string | void {
  const discordState = useDiscordState();
  if (!discordState) {
    return;
  }
  const scopes = ["identify", "bot", "applications.commands"].concat();
  const loginParams = new URLSearchParams({
    client_id: DISCORD_CLIENT_ID,
    permissions: DISCORD_BOT_PERMISSIONS, // Manage Roles | Send Messages
    redirect_uri: BOT_REDIRECT_URI,
    response_type: "code",
    scope: scopes.join(" "),
    state: discordState,
  });
  return `https://discord.com/api/oauth2/authorize?${loginParams.toString()}`;
}

export function useDiscordState(): string | undefined {
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
  return discordState;
}
