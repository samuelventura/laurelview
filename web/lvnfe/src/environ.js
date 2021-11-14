
const devURL = "ws://localhost:5001/ws"
const location = window.location
const locurl = new URL(location.href)
let wsproto = "ws:"
if (locurl.protocol==="https:") {
  wsproto = "wss:"
}
let prodURL = `${wsproto}//${locurl.host}${locurl.pathname}ws`

if (location.hostname === "laurelview.io") {
    prodURL = "wss://laurelview.io/ws"
}

const isDev = process.env.NODE_ENV === 'development'

let logEnabled = isDev

const wsURL = isDev ? devURL : prodURL

function log(...args) {
  if (logEnabled) {
    console.log(...args)
  }
}

function enableLog(enable) {
  logEnabled = enable
}

function href(path) {
  return process.env.PUBLIC_URL + path
}

const environ = {isDev, wsURL, enableLog, log, href}

window.enableLog = enableLog

export default environ
