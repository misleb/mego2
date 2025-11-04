// Load and run the WASM file in the browser.
// Modified from https://tinygo.org/docs/guides/webassembly/

const go = new Go(); // Defined in wasm_exec.js
const WASM_URL = 'main.wasm';
var wasm;

function createWebsocketClient() {
  const livePort = window.env.LIVE_PORT;
  if (!livePort) {
    return;
  }

  const socket = new WebSocket(`ws://localhost:${livePort}/ws`);
  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.cmd === "reload") {
      window.location.reload();
    }
  }
  socket.onopen = () => {
  }
  socket.onclose = () => {
    console.log("Connection Closed")  
    setTimeout(function () {  
      location.reload();  
    }, 2000); 
  }
  socket.onerror = (event) => {
    console.log("error: ", event);
  }
}

// Google OAuth configuration
window.googleAuthConfig = {
  clientId: window.env.GOOGLE_CLIENT_ID,
  redirectUri: window.env.GOOGLE_REDIRECT_URI,
};

// Function to initiate Google OAuth login
window.initiateGoogleAuth = function() {
  const params = new URLSearchParams({
    client_id: window.googleAuthConfig.clientId,
    redirect_uri: window.googleAuthConfig.redirectUri,
    response_type: 'code',
    scope: 'openid email profile',
    access_type: 'offline',
    prompt: 'consent',
  });
  
  const authUrl = `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
  
  // Open popup window for OAuth
  const width = 500;
  const height = 600;
  const left = (window.screen.width - width) / 2;
  const top = (window.screen.height - height) / 2;
  
  const popup = window.open(
    authUrl,
    'google-auth',
    `width=${width},height=${height},left=${left},top=${top}`
  );

  return new Promise((resolve, reject) => {
    const checkClosed = setInterval(() => {
      if (popup.closed) {
        clearInterval(checkClosed);
        reject(new Error('User closed the popup'));
      }
    }, 1000);
    
    // Listen for message from popup with auth code
    window.addEventListener('message', function onMessage(event) {
      if (event.data.type === 'GOOGLE_AUTH_CODE') {
        window.removeEventListener('message', onMessage);
        clearInterval(checkClosed);
        resolve(event.data.code);
      }
    }, { once: true });
  });
};

// A function to run after WebAssembly is instantiated.
function postInstantiate(obj) {
  wasm = obj.instance;
  go.run(wasm);
  createWebsocketClient();

  customElements.define('pushstate-anchor', HTMLPushStateAnchorElement, { extends: 'a' });
}

if ('instantiateStreaming' in WebAssembly) {
  WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(postInstantiate);
} else {
  fetch(WASM_URL).then(resp =>
    resp.arrayBuffer()
  ).then(bytes =>
    WebAssembly.instantiate(bytes, go.importObject).then(postInstantiate)
  )
}

class HTMLPushStateAnchorElement extends HTMLAnchorElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.addEventListener('click', this.pushStateAnchorEventListener, false);
  }

  disconnectedCallback() {
    this.removeEventListener('click', this.pushStateAnchorEventListener, false);
  }

  pushStateAnchorEventListener(event) {
    // open in new tab or open context menu (workaround for Firefox)
    if (event.ctrlKey || event.metaKey || event.which === 2 || event.which === 3) {
      return;
    }

    var href = this.getAttribute('href');
    if (!href) {
      return;
    }

    // don't pushState if the URL is for a different host
    if (href.indexOf('http') === 0 && window.location.host !== new URL(href).host) {
      return;
    }

    if (href !== window.location.pathname) {
      // push state into the history stack if the current path is not the same as the new path
      window.history.pushState(JSON.parse(this.getAttribute('state')) || window.history.state, this.getAttribute('title'), href);
      window.document.getElementById('root').scrollTo(0, 0);
    }

    // dispatch a popstate event
    try {
      var popstateEvent = new PopStateEvent('popstate', {
        bubbles: false,
        cancelable: false,
        state: window.history.state
      });

      if ('dispatchEvent_' in window) {
        // FireFox with polyfill
        window.dispatchEvent_(popstateEvent);
      } else {
        // normal
        window.dispatchEvent(popstateEvent);
      }
    } catch (error) {
      // Internet Explorer
      var evt = document.createEvent('CustomEvent');
      evt.initCustomEvent('popstate', false, false, { state: window.history.state });
      window.dispatchEvent(evt);
    }

    // prevent the default link click
    event.preventDefault();
  }
}