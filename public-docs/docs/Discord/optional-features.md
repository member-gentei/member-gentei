---
title: Optional Features
---

Please reach out to use any of these features! These require additional information to configure because they're enabled on an opt-in basis.

# Audit log channel

As a more information-rich alternative to the Discord server audit log, **gentei-bouncer#9835** can post role changes it makes to a configured channel.

![audit log message preview](/assets/bot-audit-log-message.png)

For convenient reference by a moderation team, each log message comes with a user's @tagged name, their Discord user ID, and avatar at the time of role change action.

# Multi-channel / multi-role support

Gentei supports checking multiple YouTube channels per server, mapping each membership to a separate role. When multiple membership checks are configured for a server, **gentei-bouncer#9835**'s membership check replies will list out all verified YouTube channels.

![multiple role preview](/assets/mg-check-multi.png)

If you have an audit log channel configured, log messages will be created for each role change that takes place. Each log message will mention the relevant YouTube channel.

![audit log message preview](/assets/bot-audit-log-message-multi.png)
