---
title: Bot Overview
---

Gentei's batteries-included Discord bot is called **gentei-bouncer#9835**. It's a fully automated wrapper for managing a members-only role.

Upon joining an enrolled Discord server, the bot will load memberships for the corresponding YouTube channel and start responding to the `!mg` command prefix to conduct member management.

The code for the bot is available on GitHub: [https://github.com/member-gentei/member-gentei/tree/master/bot](https://github.com/member-gentei/member-gentei/tree/master/bot)

# Server configuration

Configuration happens as part of a manual enrollment process, which happens generally over Discord or Email. See the [server onboarding guide](server-onboarding) for details and requirements.

# Commands

The screenshots below may contain outdated text and graphics, but the bot should always work as described below.

## !mg check

`!mg check` allows users to trigger a membership check on-demand, which usually completes in mere seconds. This check is rate limited to mitigate abuse.

![mg check](/assets/mg-check.png)

# Notes on operation + logging

The Discord bot is hooked up to [Google Cloud Operations](https://cloud.google.com/products/operations) for general logging and monitoring. Logs are kept for error alerting and auditing purposes and can be provided to Discord server administrators on request.

Logs of Gentei's actions can also be retrieved via the Discord server audit log. Due to technical limitations, the bot cannot currently provide detailed audit reasons in the Discord server audit log - but they do exist in Google Cloud Operations logs.

![](/assets/bot-audit-log.png)
