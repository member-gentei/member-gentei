import React from "react";
import signInDarkFocusWeb from "../../../src/assets/img/btn_google_signin_dark_focus_web.png";
import signInDarkNormalWeb from "../../../src/assets/img/btn_google_signin_dark_normal_web.png";
import signInDarkPressedWeb from "../../../src/assets/img/btn_google_signin_dark_pressed_web.png";
import { useYouTubeLoginURL } from "../LoginURL";
import { Card, CardContent, Typography } from "@mui/joy";
import { Backdrop } from "@mui/material";
import { useUser } from "../../stores/UserStore";

export function YouTubeLoginOverlay() {
  const [store] = useUser();
  return (
    <Backdrop open sx={{ position: "absolute" }}>
      <Card>
        <CardContent sx={{ textAlign: "center" }}>
          <Typography>
            Please connect your YouTube account to verify memberships.
          </Typography>
          <Typography>
            ({!store.user ? "after logging in," : null} it's the "Sign in with
            Google" button below)
          </Typography>
        </CardContent>
      </Card>
    </Backdrop>
  );
}

interface YouTubeLoginProps {
  className?: string;
}

export default function YouTubeLogin({ className }: YouTubeLoginProps) {
  const loginURL = useYouTubeLoginURL();
  const containingClassName =
    className === undefined ? "has-text-centered" : className;
  if (!loginURL) {
    return (
      <div className={containingClassName}>
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
    <div className={containingClassName}>
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
  );
}
