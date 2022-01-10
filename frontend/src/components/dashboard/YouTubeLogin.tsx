import React from "react";
import signInDarkFocusWeb from "../../../src/assets/img/btn_google_signin_dark_focus_web.png";
import signInDarkNormalWeb from "../../../src/assets/img/btn_google_signin_dark_normal_web.png";
import signInDarkPressedWeb from "../../../src/assets/img/btn_google_signin_dark_pressed_web.png";
import { useYouTubeLoginURL } from "../LoginURL";

export default function YouTubeLogin() {
  const loginURL = useYouTubeLoginURL();
  if (!loginURL) {
    return (
      <div className="has-text-centered">
        <span className="spinner"></span>
      </div>
    );
  }
  const mouseOver: React.MouseEventHandler<HTMLImageElement> = (e) => {
    e.currentTarget.setAttribute("src", signInDarkFocusWeb);
  };
  const mouseOut: React.MouseEventHandler<HTMLImageElement> = (e) => {
    e.currentTarget.setAttribute("src", signInDarkNormalWeb);
  };
  const mouseDown: React.MouseEventHandler<HTMLImageElement> = (e) => {
    e.currentTarget.setAttribute("src", signInDarkPressedWeb);
  };
  const mouseUp: React.MouseEventHandler<HTMLImageElement> = (e) => {
    e.currentTarget.setAttribute("src", signInDarkFocusWeb);
  };
  return (
    <div className="overlay is-flex is-justify-content-center is-align-content-center is-align-items-center">
      <div className="card">
        <div className="card-content">
          <p className="mb-4">
            Please connect your YouTube account below to verify memberships.
          </p>
          <div className="has-text-centered">
            <a href={loginURL}>
              <img
                src={signInDarkNormalWeb}
                onMouseOver={mouseOver}
                onMouseOut={mouseOut}
                onMouseDown={mouseDown}
                onMouseUp={mouseUp}
                alt="Sign in with Google button"
              />
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}