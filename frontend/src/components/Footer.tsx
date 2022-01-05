import React from "react";

import developedWithYouTube from "../assets/img/developed-with-youtube-sentence-case-dark.png";

interface FooterProps {
  withYouTubeImage?: boolean;
}

function Footer({ withYouTubeImage }: FooterProps) {
  let ytElement = null;
  if (withYouTubeImage) {
    ytElement = (
      <div className="is-centered developed-with-youtube">
        <img src={developedWithYouTube} alt="developed with YouTube" />
      </div>
    );
  }
  return (
    <footer className="footer pb-4">
      <div className="container">
        <div className="content has-text-centered is-size-7">
          <a href="/privacy">Privacy Policy</a> | Gentei / 限定 <br />
          Some graphics courtesy of{" "}
          <a href="https://www.irasutoya.com">いらすとや</a>
          {ytElement}
        </div>
      </div>
    </footer>
  );
}
export default Footer;
