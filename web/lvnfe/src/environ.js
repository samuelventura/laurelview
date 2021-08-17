
const devWcURL = "ws://localhost:5001/ws"
const devRtURL = "ws://localhost:5002/ws"
const location = window.location
const port = parseInt(location.port)
const prodWcURL = "ws://" + location.host + "/ws"
const prodRtURL = "ws://" + location.hostname + ":" + (port+1) + "/ws"

const isDev = process.env.NODE_ENV === 'development'

let logEnabled = isDev

const wsWcURL = isDev ? devWcURL : prodWcURL
const wsRtURL = isDev ? devRtURL : prodRtURL

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

const environ = {isDev, wsWcURL, wsRtURL, enableLog, log, href}

window.enableLog = enableLog

export default environ
