import { Action, createHook, createStore } from "react-sweet-state";
import {
  API_BASE_URL,
  authedFetchJSON,
  DISCORD_BOT_PERMISSIONS,
  LoadState,
} from "../lib/lib";

interface State {
  submit: LoadState;
  submitError?: string;
}

const initialState: State = {
  submit: LoadState.NotStarted,
};

const actions = {
  verifySubmit:
    (search: URLSearchParams, expectedState?: string): Action<State> =>
    async ({ getState, setState }) => {
      if (getState().submit >= LoadState.Started || !expectedState) {
        return;
      }
      setState({
        submit: LoadState.Started,
      });
      // check things in descending order of FUBAR -> wrong
      // ensure that the state parameter checks out
      if (search.get("state") !== expectedState) {
        setState({
          submitError: `The "state" parameter does not match what your web browser generated and has been rejected for safety and security concerns.`,
          submit: LoadState.Failed,
        });
      }
      // ensure that we have the right bot permissions
      if (search.get("permissions") !== DISCORD_BOT_PERMISSIONS) {
        setState({
          submitError: `The bot was not granted both "Manage Roles" and "Send Messages" permissions. Both are required for the bot to work properly - please reinvite the bot while allowing those permissions.`,
          submit: LoadState.Failed,
        });
      }
      // (the element that calls this action gates the call on `code`, so we can just get it here)
      // TODO: check that we have a guild_id? we don't ever use it though.
      const response = await authedFetchJSON(
        `${API_BASE_URL}/enroll-guild`,
        "POST",
        {
          code: search.get("code")!,
          permissions: search.get("permissions"),
        }
      );
      if (!response.ok) {
        const errData: {
          error: string;
          error_description?: string;
          message?: string;
        } = await response.json();
        let explainer;
        switch (errData.error_description) {
          case 'Invalid "code" in request.':
            explainer =
              "Access code is expired or otherwide invalid - please try again.";
            break;
          default:
            explainer = `${
              errData.error_description || errData.error || errData.message
            }`;
            break;
        }
        setState({
          submitError: explainer,
          submit: LoadState.Failed,
        });
        return;
      }
      // it worked - write new search and set state
      const serverData: { ID: string } = await response.json();
      const newParams = new URLSearchParams({
        server: serverData.ID,
      });
      window.location.search = `?${newParams.toString()}`;
      setState({
        submit: LoadState.Succeeded,
        submitError: undefined,
      });
    },
};

const store = createStore({ initialState, actions });

export const useEnroll = createHook(store);
