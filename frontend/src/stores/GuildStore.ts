import {
  Action,
  createContainer,
  createHook,
  createStore,
} from "react-sweet-state";
import { API_BASE_URL, authedFetchJSON, LoadState } from "../lib/lib";

interface Guild {
  ID: string;
  Name: string;
  Icon: string;
  TalentIDs: string[];
  AdminIDs: string[];
  Settings?: GuildSettings;
}

interface GuildSettings {
  RoleMapping: {
    [channelID: string]:
      | {
          ID: string;
          Name: string;
        }
      | undefined;
  };
}

interface State {
  guild?: Guild;
  guildError?: string;
  guildState: LoadState;
  saveTalentsState: LoadState;
}

const initialState: State = {
  guildState: LoadState.NotStarted,
  saveTalentsState: LoadState.NotStarted,
};

const actions = {
  load:
    (id: string, reload?: boolean): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().guildState === LoadState.Started) {
        return;
      }
      if (getState().guildState > LoadState.Started && !reload) {
        return;
      }
      setState({ guildState: LoadState.Started });
      const response = await authedFetchJSON(`${API_BASE_URL}/guild/${id}`);
      if (!response.ok) {
        if (response.status === 404) {
          setState({
            guild: undefined,
            guildError: `Discord server by ID ${id} not found`,
            guildState: LoadState.Failed,
          });
          return;
        }
        const data: {
          error: string;
          message: string;
        } = await response.json();
        setState({
          guildError: data.error || data.message,
          guildState: LoadState.Failed,
        });
        return;
      }
      setState({
        guild: await response.json(),
        guildError: undefined,
        guildState: LoadState.Succeeded,
      });
    },
  saveTalentChannels:
    (id: string, channelIDs: string[]): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().saveTalentsState === LoadState.Started) {
        return;
      }
      setState({ saveTalentsState: LoadState.Started });
      const response = await authedFetchJSON(
        `${API_BASE_URL}/guild/${id}`,
        "PATCH",
        {
          talents: channelIDs,
        }
      );
      if (!response.ok) {
        console.error(await response.json());
        setState({
          saveTalentsState: LoadState.Failed,
        });
      }
      setState({
        guild: await response.json(),
        saveTalentsState: LoadState.Succeeded,
      });
    },
};

const store = createStore({ initialState, actions });

export const GuildContainer = createContainer(store);
export const useGuild = createHook(store);
