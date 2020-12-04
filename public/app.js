document.addEventListener('DOMContentLoaded', function () {
    var firestore = firebase.firestore()
    const appQuerySelector = "#app"
    Vue.component("yt-info", {
        template: "#yt-info-template",
        props: {
            dcyoutubeaccount: {
                type: Function,
                default: () => { () => { } },
            },
            "channelid": {
                default: () => { return "" },
            },
            "channels": {
                type: Array,
                default: () => { return [] },
            },
            "memberships": {
                type: Array,
                default: () => { return [] },
            },
            "loginUrl": {
                type: String,
                required: true,
            }
        },
        data: function () {
            return {
                deleting: false,
            }
        },
        computed: {
            channelHref: function () {
                if (this.channelid != "") {
                    return "https://youtube.com/channel/" + this.channelid;
                }
                // joke: "oki doki boomer" autoplay video. If this ever displays, it's a bug.
                return "https://youtube.com/channel/UCTUHzVzRwN_2x13IWQ9QVNg"
            },
            memberMap: function () {
                let memberMap = {}
                if (!this.memberships) {
                    return memberMap
                }
                this.memberships.forEach(element => {
                    memberMap[element.path] = true
                })
                return memberMap
            },
            channelMemberships: function () {
                var memberOf = [],
                    nonMember = [],
                    memberMap = this.memberMap
                if (!this.channels) {
                    return []
                }
                this.channels.forEach(element => {
                    let relevantArray
                    if (!!memberMap[element.path]) {
                        memberOf.push({
                            path: element.path,
                            docRef: element,
                            isMember: true,
                        })
                    } else {
                        nonMember.push({
                            path: element.path,
                            docRef: element,
                            isMember: false,
                        })
                    }
                })
                return memberOf.concat(nonMember)
            },
        },
        methods: {
            disconnect: async function () {
                this.deleting = true
                await this.dcyoutubeaccount()
            }
        }
    })
    Vue.component("yt-channel-card", {
        template: "#yt-channel-card-template",
        props: ["docref", "isMember"],
        data: function () {
            return {
                doc: null,
            }
        },
        computed: {
            channelHref: function () {
                if (!!this.doc) {
                    return "https://youtube.com/channel/" + this.doc.ChannelID
                }
                return ""
            }
        },
        mounted: async function () {
            if (!this.docref) {
                return null
            }
            let doc = await this.docref.get()
            this.doc = doc.data()
        }
    })
    Vue.component("login", {
        template: "#login-template",
        props: {
            loginType: String,
            setLoginError: Function,
            setCurrentRoute: Function,
        },
        methods: {
            xhrSuccess: async function(xhr) {
                await firebase.auth().signInWithCustomToken(xhr.response["jwt"])
                this.setCurrentRoute("/app")
            },
            xhrFail: function(xhr) {
                console.log("fail")
                if (xhr.responseText === "") {
                    this.setLoginError({"message": "Unexpected error signing in with Discord - probably a browser CORS error?"})
                } else if (!!xhr.response) {
                    this.setLoginError(xhr.response)
                } else {
                    this.setLoginError({"message": "unhandled error - " + xhr.responseText})
                }
                this.setCurrentRoute("/app")
            }
        },
        mounted: function() {
            var that = this
            let loginURL = "https://us-central1-member-gentei.cloudfunctions.net/Auth?service=" + this.loginType
            let locationParams = new URL(document.location).searchParams
            if (!locationParams.get("code")) {
                console.log("no OAuth code found")
                return
            }
            let xhr = new XMLHttpRequest();
            xhr.open("POST", loginURL, true);
            xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
            xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
            xhr.setRequestHeader('Accept', 'application/json');
            xhr.onreadystatechange = function () {
                if (xhr.readyState > 3) {
                    if (xhr.status == 200) {
                        that.xhrSuccess(xhr)
                    } else {
                        that.xhrFail(xhr)
                    }
                }
            }
            xhr.send(locationParams)
        },
    })
    const app = new Vue({
        "el": appQuerySelector,
        data: {
            "currentRoute": window.location.pathname,
            "user": null,
            "loginError": {},
            "userData": {},
            "loaded": false,
        },
        computed: {
            loginURLs: function () {
                let loginURL = encodeURIComponent(window.location.protocol + "//" + window.location.host + "/login/discord")
                return {
                    discord: "https://discord.com/api/oauth2/authorize?client_id=768486576388177950&redirect_uri=" + loginURL + "&response_type=code&scope=identify%20guilds",
                    youtube: "https://accounts.google.com/o/oauth2/v2/auth?client_id=649732146530-s4cj4tqo2impojg7ljol2chsuj1us81s.apps.googleusercontent.com&redirect_uri=https%3A%2F%2Fmember-gentei.tindabox.net%2Flogin%2Fyoutube&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fyoutube.force-ssl&access_type=offline&prompt=consent",
                }
            },
        },
        watch: {
            user: function (newUser, oldUser) {
                if (oldUser === null && newUser !== null) {
                    // start a realtime update thing
                    var that = this
                    firestore.collection("users").doc(that.user.uid).onSnapshot(doc => {
                        that.userData = doc.data();
                    })
                }
            },
        },
        methods: {
            dcYouTubeAccount: async function () {
                await firestore.collection("users").doc(this.user.uid).collection("private").doc("youtube").delete()
            },
            setLoginError: function(loginError) {
                this.loginError = loginError
            },
            setCurrentRoute: function(newRoute) {
                this.currentRoute = newRoute
                history.pushState({}, "", document.location.origin + "/app")
            },
            logout: function () {
                firebase.auth().signOut()
                this.user = null
            }
        },
        mounted: function () {
            var that = this
            firebase.auth().onAuthStateChanged(user => {
                if (user) {
                    that.user = user
                } else {
                    that.user = null
                    that.loaded = true
                }
            })
            this.$nextTick(function () {
                document.querySelector(appQuerySelector).classList.remove("is-hidden")
            })
        },
    })
    window.addEventListener("popstate", () => {
        app.currentRoute = window.location.pathname
    })
})
