import React, { useEffect, useReducer } from 'react'

import Modal from 'react-bootstrap/Modal';

import ItemBrowser from "./app/ItemBrowser"
import "./App.css"

import socket from "./socket"
import env from "./environ"

function App() {
  
  function reducer(state, {name, args, session}) {
    // called twice on purpose by reactjs 
    // to detect side effects on strict mode
    // reducer must be pure
    switch(name){
      case "all": {
        const next = {...state}
        next.session = session
        next.online = true
        next.items = {}
        args.items.forEach(item => { 
          next.items[item.id] = item
        })
        return next
      }
      case "create": {
        const next = {...state}
        next.items[args.id] = args
        return next
      }
      case "delete": {
        const next = {...state}
        delete next.items[args.id]
        return next
      }
      case "update": {
        const next = {...state}
        next.items[args.id].name = args.name
        next.items[args.id].json = args.json
        return next
      }
      case "close": {
        //flickers on navigating back (reconnect)
        const next = {...state}
        next.items = {}
        next.online = false
        next.session = null
        return next
      }
      case "send": {
        const next = {...state}
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

  function handleDispatch({name, args}) {
    switch(name) {
      case "create":
      case "delete":
      case "update":
        state.send({name, args})
        break
      default:
        env.log("Unknown mutation", name, args)
    }
  }

  useEffect(() => {
    return socket.createWc(dispatch, "/index")
  }, [])

  return (
    <div className="App">
      <ItemBrowser 
        state={state} 
        dispatch={handleDispatch} />
      <Modal
        show={!state.online}
        backdrop="static"
        keyboard={false}
        centered
      >
        <Modal.Header>
          <Modal.Title>Connecting</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          Connecting to backend...
        </Modal.Body>
      </Modal>
    </div>
  )
}

export default App
