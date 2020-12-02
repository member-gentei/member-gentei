---
title: YouTube Permissions
---

One of the common comments about Gentei are that the permissions look kinda sketchy.

![oauth-permission](/assets/oauth-permissions.png)

The short response to that is this is unavoidable for the way that Gentei currently verifies memberships. I wish this could work with a less restrictive permission set, but I can at least explain why Gentei needs it.

# YouTube API permissions

Gentei depends on reading out comment threads on a members-only video, and [the associated YouTube Data API endpoint](https://developers.google.com/youtube/v3/docs/commentThreads/list) requires a specific OAuth permission scope called `https://www.googleapis.com/auth/youtube.force-ssl`.

![youtube-force-ssl](/assets/youtube-force-ssl.png)

Referencing that scope with YouTube's [table of available OAuth scopes](https://developers.google.com/youtube/v3/guides/auth/server-side-web-apps#identify-access-scopes), this maps to the "See, edit, and permanently delete your YouTube videos, ratings, comments and captions" message that users see in the top screenshot above.

There's not really anything we can do about that past asking YouTube to reclassify this API endpoint under a separate OAuth scope and hope that it happens. ¯\\\_(ツ)_/¯
