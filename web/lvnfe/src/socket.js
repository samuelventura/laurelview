import env from "./environ"

function create(dispatch, path) {
  let toms = 0
  let to = null
  let ws = null
  let closed = true
  let disposed = false

  function safe(action) {
    try { action() }
    catch(e) { env.log("exception", e) }
  }

  function dispose() {
    env.log("dispose", disposed, closed, to, ws)
    disposed = true
    if (to) clearTimeout(to)
    if (ws) safe(() => ws.close())
  }

  function send(msg) {
    env.log("ws.send", disposed, closed, msg)
    if (disposed) return
    if (closed) return
    safe(() => ws.send(JSON.stringify(msg)) )
  }

  function connect() {
    //immediate error when navigating back
    //toms is workaround for trottled reconnection
    //safari only, chrome and firefox work ok
    let url = env.wsURL + path
    ws = new WebSocket(url)
    env.log("connect", to, url, ws)
    ws.onclose = (event) => {  
      env.log("ws.close", event)
      closed = true
      if (disposed) return
      dispatch({name: "close"})
      to = setTimeout(connect, toms)
      toms += 1000 
      toms %= 4000
    }
    ws.onmessage = (event) => {
      //env.log("ws.message", event, event.data)
      const msg = JSON.parse(event.data)
      env.log("ws.message", msg)
      dispatch(msg)
    }
    ws.onerror = (event) => {
      env.log("ws.error", event)
    }
    ws.onopen = (event) => {
      env.log("ws.open", event)
      dispatch({name: "send", args: send})
      closed = false
      toms = 0
    }
  }
  to = setTimeout(connect, 0)
  return dispose
}

function send(msg) {
  env.log("nop.send", msg)
}

var socket = {create, send}

export default socket
