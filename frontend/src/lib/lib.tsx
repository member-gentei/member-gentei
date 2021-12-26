/** No trailing slash! */
export const API_BASE_URL = process.env.REACT_APP_PROD
  ? "https://gentei-api.tindabox.net"
  : "http://localhost:5000";

export const DISCORD_CLIENT_ID = process.env.REACT_APP_PROD
  ? "768486576388177950"
  : "924507400139071528";

export const REDIRECT_URI = process.env.REACT_APP_PROD
  ? "https://gentei.tindabox.net/login/discord"
  : "http://localhost:3000/login/discord";

/** General-purpose load state enum. */
export enum LoadState {
  NotStarted,
  Started,
  Loaded,
  Succeeded,
  Failed,
}

export function authedFetchJSON(
  url: string,
  method = "GET",
  body?: { [key: string]: any }
): ReturnType<typeof fetch> {
  return fetch(url, {
    method,
    headers: {
      "content-type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(body),
  });
}
