---
title: Server Onboarding
---

After some prep work, **gentei-bouncer#9835** can be added to an enrolled Discord fan server.

# Requirements

To conduct membership administration, the bot currently expects to:

* Manage a members-only Discord role
* Respond to membership check requests **in any channel it can read and respond to**

Most administrators will want to restrict the channels the bot can read via various role restrictions.

# Enrollment process

[Sign up to enroll a community](https://forms.gle/rr4Psqbzz1Nhuqno6) and I'll get back to you in a day or two. In the meantime, feel free to read up on the below! 

## Collecting integration info

After verifying the request, I'll be reaching out - if via Discord, as `のヮの＃6969` - to collect some information about the Discord server. I'll be asking for:

* The Discord Guild ID, or simply an invite to the appropriate Discord server
* The members-only role ID that the bot will be managing
* How the Discord server's admin/mod team would like to receive Gentei announcements or updates
* Preferred contact info for the Discord server's admin/mod team

If you need help with any of the above, let me know!

## Invite the Discord bot

You can join **gentei-bouncer#9835** to your Discord server of choice by clicking the following enrollment link: [https://discord.com/oauth2/authorize?client_id=768486576388177950&permissions=268438528&scope=bot](https://discord.com/oauth2/authorize?client_id=768486576388177950&permissions=268438528&scope=bot).

## Configure the bot's Discord role

Upon joining, **gentei-bouncer** will have a role called `Gentei`, and by default Discord will put this at the bottom of the role list, right above `@everyone`. You can rename this role if you want.

![roles on join](/assets/roles-on-join.png)

To allow **gentei-bouncer** to manage the members-only role, drag its `Gentei` role above that role and save the change.

![role in correct position](/assets/dragged-role.png)

To learn more about why role order matters, see Discord's [Role Management 101 help article](https://support.discord.com/hc/en-us/articles/214836687-Role-Management-101).

### Recommended: restrict the bot's channels

To prevent **gentei-bouncer** from responding to `!mg check` in random channels, we recommend following the common permission scheme of granting read + send permissions to a special purpose channel for interacting with the bot. Testers have named this channel something like `#member-gate` or `#member-check`.

Most Discord servers with a "I have read the rules" role already have a permissions structure that easily accommodates this permission scheme, as **gentei-bouncer** will not inherit potentially disruptive permissions from the special `@everyone` role.

## Co-existing with other role assignment methods

**gentei-bouncer** enforces membership verification both by automatically adding verified users to the members-only role and by automatically removing users from the role whose memberships cannot be verified. 

This will likely run afoul of alternate membership verification methods that a Discord server mod team might want to use, so we know of a couple ways to address that:

* Run alternate methods on an equivalent members-only role that **gentei-bouncer** will not manage.
* Create/integrate your own bot with the [Gentei API](https://mark-ignacio.github.io/member-gentei/api.html).

If neither of these are acceptable options, please reach out and we can talk about how to handle things.

## Optional features

To make administrative life better - but by adding some complexity - [optional features](optional-features) can be configured for the server!

# Offboarding

To remove Gentei from a Discord server, kick the bot! It can be re-invited at any time to continue managing its configured member role, but kicking the bot means that all configuration done to its role will be lost.