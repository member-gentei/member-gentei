import React from 'react'

function ShutdownModal(): JSX.Element {
    return (
        <div className="modal is-active">
            <div className="modal-background"></div>
            <div className="modal-content">
                <div className="message is-d    ark">
                    <div className="message-body">
                        <div className="content">
                            <p>Gentei has been shut down on May 23, 2021 due to YouTube API changes that no longer allow membership verification through access control on the <a href="https://developers.google.com/youtube/v3/docs/commentThreads/list">CommentThreads.list</a> API call.</p>
                            <p>Gentei's last actually-working day of May 3 had 11,381 users with 33,963 memberships across 62 Discord servers dedicated to 57 VTubers and 1 actual human being. Hope you all appreciated it while it lasted!</p>
                            <p>Feel free to hop into <code>#gentei-限定</code> on the Hololive Creators Club Discord server for discussions of alternate membership verification automation to succeed the Gentei bot and app.</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default ShutdownModal