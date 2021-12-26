import { Action, createHook, createStore } from "react-sweet-state";
import { API_BASE_URL, authedFetchJSON, LoadState } from "../lib/lib";

interface State {
  user?: {
    ID: string;
    FullName: string;
    AvatarHash: string;
    YouTube: {
      ID: string;
      Valid: boolean;
    };
    Memberships: {};
  };
  userLoad: LoadState;
  discordLogin: LoadState;
  discordLoginError?: { [key: string]: any };
}

const initialState: State = {
  userLoad: LoadState.NotStarted,
  discordLogin: LoadState.NotStarted,
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
      if (response.status === 400) {
        const data: { message: string } = await response.json();
        if (data.message === "missing or malformed jwt") {
          setState({
            user: undefined,
            userLoad: LoadState.Succeeded,
          });
          return;
        }
      }
      setState({
        user: await response.json(),
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
        { code, state }
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
  logout:
    (): Action<State> =>
    async ({ getState, setState }) => {
      if (!getState().user) {
        return;
      }
      setState({ user: undefined });
      await fetch(`${API_BASE_URL}/logout`);
    },
};

const store = createStore({ initialState, actions });

export const useUser = createHook(store);
