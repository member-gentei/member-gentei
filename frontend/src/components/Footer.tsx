import React from "react";
import { RiGithubFill } from "react-icons/ri";

import developedWithYouTube from "../assets/img/developed-with-youtube-sentence-case-dark.png";
import { Container, Link, Typography } from "@mui/joy";

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
          <Link href="https://github.com/member-gentei/member-gentei">
            <RiGithubFill />
          </Link>
          <br />
          <Link href="/tos">Terms of Service</Link> |{" "}
          <Link href="/privacy">Privacy Policy</Link> | Gentei / 限定 <br />
          Some graphics courtesy of{" "}
          <Link href="https://www.irasutoya.com">いらすとや</Link>
        </Typography>
        {ytElement}
      </Container>
    </footer>
  );
}
export default Footer;
