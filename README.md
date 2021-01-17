# member-gentei

_Gentei_ is a Firebase-and-friends application that monitors YouTube channel memberships via client credentials. It facilitates automatic, members-only role assignment for Discord fan server administrators.

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