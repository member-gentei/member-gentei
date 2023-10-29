import React from "react";
import { RiGithubFill } from "react-icons/ri";

import developedWithYouTube from "../assets/img/developed-with-youtube-sentence-case-dark.png";
import { Container, Typography } from "@mui/joy";

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
    <footer className="footer pt-6 pb-4">
      <Container sx={{ textAlign: "center" }}>
        <Typography level="body-sm">
          <div>
            <a href="https://github.com/member-gentei/member-gentei">
              <RiGithubFill />
            </a>
            <br />
            <a href="/tos">Terms of Service</a> |{" "}
            <a href="/privacy">Privacy Policy</a> | Gentei / 限定 <br />
            Some graphics courtesy of{" "}
            <a href="https://www.irasutoya.com">いらすとや</a>
            {ytElement}
          </div>
        </Typography>
      </Container>
    </footer>
  );
}
export default Footer;
