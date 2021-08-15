
const devWcURL = "ws://localhost:5001/ws"
const devRtURL = "ws://localhost:5002/ws"
const prodURL = process.env.PUBLIC_URL + "/ws"

const isDev = process.env.NODE_ENV === 'development'

let logEnabled = isDev

const wsWcURL = isDev ? devWcURL : prodURL
const wsRtURL = isDev ? devRtURL : prodURL

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

export default environ
