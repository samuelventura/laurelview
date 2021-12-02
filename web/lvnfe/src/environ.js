
const locurl = new URL(window.location.href)
let wsproto = "ws:"
if (locurl.protocol === "https:") {
  wsproto = "wss:"
}
let wsURL = `${wsproto}//${locurl.host}${locurl.pathname}ws`

const isDev = process.env.NODE_ENV === 'development'

let logEnabled = isDev

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

const environ = { isDev, wsURL, enableLog, log, href }

window.enableLog = enableLog

export default environ
