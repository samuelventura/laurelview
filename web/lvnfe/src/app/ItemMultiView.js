import React, { useReducer, useEffect } from 'react'

import Modal from 'react-bootstrap/Modal';

import ItemDisplay from "./ItemDisplay"

import socket from "../socket"
import env from "../environ"

function ItemMultiView(props) {

  function reducer(state, { name, args, session }) {
    switch (name) {
      case "query": {
        const next = { ...state }
        next.queries[args.index] = args
        return next
      }
      case "close": {
        const next = { ...state }
        next.send = socket.send
        return next
      }
      case "send": {
        const next = { ...state }
        next.queries = props.items.map(() => { return {}})
        next.send = args
        return next
      }
      default:
        env.log("Unknown mutation", name, args, session)
        return state
    }
  }

  const initial = {
    send: socket.send,
    queries: []
  }

  const [state, dispatch] = useReducer(reducer, initial)

  useEffect(() => {
    function handler({ name, args, session }) {
      dispatch({ name, args, session })
      switch (name) {
        case "send": {
          const name = "setup"
          const items = props.items.map(item => JSON.parse(item.json))
          env.log({ name, args: { items } })
          args({ name, args: { items } })
          break
        }
        default: //linter complains
      }
    }
    return socket.createRt(handler, "/index")
  }, [props])

  function handleHide() {
    const action = "cancel"
    props.handler({ action })
  }

  function query(index) {
    let q = state.queries[index]
    q = q || {}
    q.index = index
  }

  const displays = props.items.map((item, i) => {
    return <div key={item.id}>
        <h3>{item.name}</h3>
        <ItemDisplay query={query(i)} />
    </div>
  })

  return (
    <Modal show={props.show} onHide={handleHide} backdrop="static" centered>
      <Modal.Header closeButton>
        <Modal.Title>Multi View</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {displays}
      </Modal.Body>
    </Modal>
  )
}

export default ItemMultiView
