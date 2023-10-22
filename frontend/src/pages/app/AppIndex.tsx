import classNames from "classnames";
import React, { Fragment, useRef, useState } from "react";
import { RiArrowDropDownLine } from "react-icons/ri";
import { SiDiscord } from "react-icons/si";
import { Link, Outlet } from "react-router-dom";
import logo128 from "../../assets/img/logo-128.png";
import Footer from "../../components/Footer";
import { useDiscordLoginURL } from "../../components/LoginURL";
import { LoadState, useWindowSize } from "../../lib/lib";
import { useUser } from "../../stores/UserStore";
import {
  AppBar,
  IconButton,
  Toolbar,
  Menu,
  MenuItem,
  List,
  ListItem,
} from "@mui/material";
import { Avatar, Box, Button, Container, Tooltip, Typography } from "@mui/joy";

function AppIndex() {
  const actions = useUser()[1];
  actions.getMe();
  return (
    <Fragment>
      <AppBar component="nav" position="static">
        <Toolbar disableGutters>
          <IconButton>
            <Link to="/app" className="navbar-item">
              <img src={logo128} alt="Gentei bot logo" height={40} />
            </Link>
          </IconButton>
          <Box sx={{ flexGrow: 1 }} />
          <Box sx={{ flexGrow: 0 }}>
            <AuthButtons />
          </Box>
        </Toolbar>
      </AppBar>
      <Container>
        <Box sx={{ mt: 4 }}>
          <Outlet />
        </Box>
      </Container>
      <Footer withYouTubeImage />
    </Fragment>
  );
}

function AuthButtons() {
  const [store, actions] = useUser();
  const [menuActive, setMenuActive] = useState(false);
  const loginURL = useDiscordLoginURL();
  const iconButtonRef = useRef(null);
  let innards = null;
  let menu = null;
  if (store.userLoad === LoadState.Succeeded && !!store.user) {
    const user = store.user!;
    const avatarURL = `https://cdn.discordapp.com/avatars/${user.ID}/${user.AvatarHash}.webp?size=64`;
    const logout: React.MouseEventHandler<HTMLAnchorElement> = (e) => {
      e.preventDefault();
      actions.logout();
    };
    innards = (
      <Avatar alt={`Discord avatar for ${user.FullName}`} src={avatarURL} />
    );
    menu = (
      <Menu
        anchorEl={iconButtonRef.current}
        open={menuActive}
        onClose={() => setMenuActive(false)}
        keepMounted
      >
        <MenuItem>
          <Typography textAlign="center">{user.FullName}</Typography>
        </MenuItem>
        <MenuItem>
          <Typography textAlign="center">
            <Link to="/logout" onClick={logout}>
              Sign out
            </Link>
          </Typography>
        </MenuItem>
      </Menu>
    );
  } else {
    return (
      <List>
        <ListItem>
          <Button
            component="a"
            href={loginURL || "#"}
            startDecorator={<SiDiscord />}
            variant="soft"
          >
            Register / Sign in with Discord
          </Button>
        </ListItem>
      </List>
    );
  }
  return (
    <Fragment>
      <Tooltip title="Open settings">
        <IconButton
          ref={iconButtonRef}
          onClick={() => setMenuActive((v) => !v)}
        >
          {innards}
        </IconButton>
      </Tooltip>
      {menu}
    </Fragment>
  );
}

export default AppIndex;
