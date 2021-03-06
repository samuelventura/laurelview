import React, { useEffect, useReducer } from 'react'

import Navbar from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';

import ItemBrowser from "./app/ItemBrowser"
import "./App.css"

import socket from "./socket"
import env from "./environ"

function App() {

  function reducer(state, { name, args, session }) {
    // called twice on purpose by reactjs 
    // to detect side effects on strict mode
    // reducer must be pure
    switch (name) {
      case "all": {
        const next = { ...state }
        next.session = session
        next.online = true
        next.items = {}
        args.forEach(item => {
          next.items[item.id] = item
        })
        return next
      }
      case "create": {
        const next = { ...state }
        next.items[args.id] = args
        return next
      }
      case "delete": {
        const next = { ...state }
        delete next.items[args]
        return next
      }
      case "update": {
        const next = { ...state }
        next.items[args.id].name = args.name
        next.items[args.id].json = args.json
        return next
      }
      case "close": {
        //flickers on navigating back (reconnect)
        const next = { ...state }
        //keep view items so multiview wont hide
        //next.items = {}
        next.online = false
        next.session = null
        next.send = socket.send
        return next
      }
      case "open": {
        const next = { ...state }
        next.send = args
        return next
      }
      default:
        env.log("Unknown mutation", name, args, session)
        return state
    }
  }

  const initial = {
    items: {},
    online: false,
    session: null,
    send: socket.send
  }

  const [state, dispatch] = useReducer(reducer, initial)

  function handleDispatch({ name, args }) {
    switch (name) {
      case "create":
      case "delete":
      case "update":
        state.send({ name, args })
        break
      default:
        env.log("Unknown mutation", name, args)
    }
  }

  useEffect(() => {
    return socket.create(dispatch, "/db")
  }, [])

  function offline() {
    if (state.online) return
    return (<Navbar bg="dark" variant="dark">
      <Container>
        <Navbar.Text>
          Connecting to backend...
        </Navbar.Text>
      </Container>
    </Navbar>)
  }

  return (
    <div className="App">
      {offline()}
      <ItemBrowser
        state={state}
        dispatch={handleDispatch} />
    </div>
  )
}

export default App
