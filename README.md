# member-gentei

_Gentei_ is a Firebase-and-friends application that monitors YouTube channel memberships via client credentials. It facilitates automatic, members-only role assignment for Discord fan server administrators.


## Shutdown!

Gentei has been shut down on May 23, 2021 due to YouTube API changes that no longer allow membership verification through access control on the [CommentThreads.list](https://developers.google.com/youtube/v3/docs/commentThreads/list) API call.

Gentei's last actually-working day of May 3 had 11,381 users with 33,963 memberships across 62 Discord servers dedicated to 57 VTubers and 1 actual human being. Hope you all appreciated it while it lasted!

Feel free to hop into `#gentei-限定` on the Hololive Creators Club server for discussions of alternate membership verification automation to succeed the Gentei bot and app.

## API

OpenAPI documentation: https://docs.member-gentei.tindabox.net/api.html

## Project layout

* `/bot` - the code for **gentei-bouncer#9835**
* `/docs` - https://docs.member-gentei.tindabox.net
* `/functions` - Google Cloud Function code
* `/init` - systemd, package files
* `/gentei` - React Typescript app for https://member-gentei.tindabox.net
* `/jobs` - systemd timer jobs
* `/pkg` - stuff for reuse across multiple functions + packages
* `/public-docs` - [Doctave](https://github.com/Doctave/doctave) template that generates `/docs`
* `/tools` - command line tooling for administrative actions
* `openapi.yml` - the authoritative OpenAPI v3 spec

## Additional credit

Gentei's Discord avatar and favicon is by [@Dakuma_Art](https://twitter.com/Dakuma_Art). Check her out!