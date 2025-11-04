window.env = {
	LIVE_PORT: window.location.hostname === "localhost" ? 38919 : null,
	GOOGLE_CLIENT_ID: "920235156207-21u7vs6ccabv14itrlp6cgapmi4ki6b3.apps.googleusercontent.com",
	GOOGLE_REDIRECT_URI: window.location.origin + '/google-callback.html',
}
