import { useEffect, useState } from "react";

/** No trailing slash! */
export const API_BASE_URL = process.env.REACT_APP_PROD
  ? "https://gentei-api.tindabox.net"
  : "http://localhost:5000";

export const DISCORD_CLIENT_ID = process.env.REACT_APP_PROD
  ? "768486576388177950"
  : "924507400139071528";

export const DISCORD_REDIRECT_URI = process.env.REACT_APP_PROD
  ? "https://gentei.tindabox.net/login/discord"
  : "http://localhost:3000/login/discord";

export const BOT_REDIRECT_URI = process.env.REACT_APP_PROD
  ? "https://gentei.tindabox.net/app/enroll"
  : "http://localhost:3000/app/enroll";

export const YOUTUBE_CLIENT_ID = process.env.REACT_APP_PROD
  ? "649732146530-s4cj4tqo2impojg7ljol2chsuj1us81s.apps.googleusercontent.com"
  : "649732146530-med3rfenvlo8ahfdcv69dkntd0f1jcj7.apps.googleusercontent.com";

export const YOUTUBE_REDIRECT_URI = process.env.REACT_APP_PROD
  ? "https://gentei.tindabox.net/login/youtube"
  : "http://localhost:3000/login/youtube";

export const DISCORD_BOT_PERMISSIONS = "268437504";

/** 0001-01-01 00:00:00+00:00, which Go uses as the zero value for time.Time. */
export const ZeroTime = -62135596800;

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

export function useWindowSize() {
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);
  useEffect(() => {
    function handler() {
      setWidth(window.innerWidth);
      setHeight(window.innerHeight);
    }
    window.addEventListener("resize", handler);
    handler();
    return () => {
      window.removeEventListener("resize", handler);
    };
  });
  return {
    width: width,
    height: height,
  };
}
