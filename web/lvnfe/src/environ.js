
const devURL = "ws://localhost:5001/ws"
const prodURL = process.env.PUBLIC_URL + "/ws"

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

export default environ
