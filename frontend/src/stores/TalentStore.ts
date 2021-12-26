import { Action, createHook, createStore } from "react-sweet-state";
import { API_BASE_URL, authedFetchJSON, LoadState } from "../lib/lib";

export interface Talent {
  ID: string;
  Name: string;
  Thumbnail: string;
}

interface State {
  talentsByID: { [id: string]: Talent | undefined };
  loadAllState: LoadState;
}

const initialState: State = {
  talentsByID: {},
  loadAllState: LoadState.NotStarted,
};

const actions = {
  loadAll:
    (reload = false): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().loadAllState >= LoadState.Started && !reload) {
        return;
      }
      setState({ loadAllState: LoadState.Started });
      const response = await authedFetchJSON(`${API_BASE_URL}/talents`);
      if (!response.ok) {
        console.error(await response.json());
        setState({
          loadAllState: LoadState.Failed,
        });
        return;
      }
      // generate talentsByID
      const talents: Talent[] = await response.json();
      let talentsByID: State["talentsByID"] = {};
      talents.forEach((t) => {
        talentsByID[t.ID] = t;
      });
      setState({
        talentsByID,
        loadAllState: LoadState.Succeeded,
      });
    },
};

const store = createStore({ initialState, actions });

export const useTalents = createHook(store);
