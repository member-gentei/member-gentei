# member-gentei

_Gentei_ is a Firebase-and-friends application that monitors YouTube channel memberships via client credentials. It facilitates automatic, members-only role assignment for Discord fan server administrators.

## API

OpenAPI documentation: https://mark-ignacio.github.io/member-gentei/api.html

## Project layout

* `/bot` - the code for **gentei-bouncer#9835**
* `/docs` - https://docs.member-gentei.tindabox.net
* `/functions` - Google Cloud Function code
* `/init` - systemd, package files
* `/jobs` - systemd timer jobs
* `/pkg` - stuff for reuse across multiple functions + packages
* `/public` - https://member-gentei.tindabox.net
* `/public-docs` - [Doctave](https://github.com/Doctave/doctave) template that generates `/docs`
* `/tools` - command line tooling for administrative actions
* `openapi.yml` - the authoritative OpenAPI v3 spec
