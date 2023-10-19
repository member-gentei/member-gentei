import { Action, createHook, createStore } from "react-sweet-state";
import { API_BASE_URL, authedFetchJSON, LoadState } from "../lib/lib";

interface State {
  user?: {
    ID: string;
    FullName: string;
    AvatarHash: string;
    LastRefreshed: number;
    YouTube: {
      ID: string;
      Valid: boolean;
    };
    Memberships?: {
      [talentID: string]:
        | {
            LastVerified: number;
            Failed: boolean;
          }
        | undefined;
    };
    ServerAdmin?: string[];
    Servers?: string[];
  };
  derived: {
    sortedServers: string[];
  };
  userLoad: LoadState;
  discordLogin: LoadState;
  discordLoginError?: { [key: string]: any };
  youtubeLogin: LoadState;
  youtubeLoginError?: { [key: string]: any };
  disconnect: LoadState;
  disconnectError?: string;
  remove: LoadState;
  removeError?: string;
}

const initialState: State = {
  userLoad: LoadState.NotStarted,
  derived: {
    sortedServers: [],
  },
  discordLogin: LoadState.NotStarted,
  youtubeLogin: LoadState.NotStarted,
  disconnect: LoadState.NotStarted,
  remove: LoadState.NotStarted,
};

const actions = {
  getMe:
    (reload = false): Action<State> =>
    async ({ getState, setState }) => {
      const loadState = getState().userLoad;
      if (loadState === LoadState.Started) {
        return;
      } else if (!reload && loadState >= LoadState.Loaded) {
        return;
      }
      setState({ userLoad: LoadState.Started });
      const response = await authedFetchJSON(`${API_BASE_URL}/me`);
      switch (response.status) {
        case 400:
          const data: { message: string } = await response.json();
          if (data.message === "missing or malformed jwt") {
            setState({
              user: undefined,
              userLoad: LoadState.Succeeded,
            });
            return;
          }
          break;
        case 401:
          setState({
            user: undefined,
            userLoad: LoadState.Succeeded,
          });
          return;
      }
      const user: State["user"] = await response.json();
      // concat servers
      const serverSet = new Set(
        (user?.Servers || []).concat(user?.ServerAdmin || []),
      );
      let sortedServers = [];
      for (const key of serverSet.keys()) {
        sortedServers.push(key);
      }
      setState({
        user: user,
        derived: {
          sortedServers: sortedServers.sort(),
        },
        userLoad: LoadState.Succeeded,
      });
    },
  loginDiscord:
    (code: string, state: string): Action<State> =>
    async ({ getState, setState }) => {
      setState({ discordLogin: LoadState.Started });
      const response = await authedFetchJSON(
        `${API_BASE_URL}/login/discord`,
        "POST",
        { code, state },
      );
      if (!response.ok) {
        setState({
          discordLogin: LoadState.Failed,
          discordLoginError: await response.json(),
        });
      } else {
        setState({
          user: await response.json(),
          discordLogin: LoadState.Succeeded,
        });
      }
    },
  loginYouTube:
    (code: string, state: string): Action<State> =>
    async ({ getState, setState }) => {
      setState({ youtubeLogin: LoadState.Started });
      const response = await authedFetchJSON(
        `${API_BASE_URL}/login/youtube`,
        "POST",
        { code, state },
      );
      if (!response.ok) {
        setState({
          youtubeLogin: LoadState.Failed,
          youtubeLoginError: await response.json(),
        });
        return;
      }
      const result: {
        ChannelID: string;
      } = await response.json();
      const user = getState().user!;
      setState({
        youtubeLogin: LoadState.Succeeded,
        user: {
          ...user,
          YouTube: {
            ID: result.ChannelID,
            Valid: true,
          },
        },
      });
      return;
    },
  logout:
    (): Action<State> =>
    async ({ getState, setState }) => {
      if (!getState().user) {
        return;
      }
      setState({ user: undefined });
      await authedFetchJSON(`${API_BASE_URL}/logout`, "POST");
    },
  disconnectYouTube:
    (): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().disconnect > LoadState.NotStarted) {
        return;
      }
      setState({ disconnect: LoadState.Started });
      const response = await authedFetchJSON(
        `${API_BASE_URL}/me/youtube`,
        "DELETE",
      );
      if (!response.ok) {
        setState({
          disconnect: LoadState.Failed,
          disconnectError: await response.json(),
        });
        return;
      }
      setState({
        disconnect: LoadState.NotStarted,
        user: await response.json(),
      });
    },
  remove:
    (): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().remove > LoadState.NotStarted) {
        return;
      }
      setState({ remove: LoadState.Started });
      const response = await authedFetchJSON(`${API_BASE_URL}/me`, "DELETE");
      if (!response.ok) {
        setState({
          remove: LoadState.Failed,
          removeError: await response.json(),
        });
        return;
      }
      setState({
        remove: LoadState.NotStarted,
        userLoad: LoadState.Succeeded,
        user: undefined,
      });
    },
};

const store = createStore({ initialState, actions });

export const useUser = createHook(store);
