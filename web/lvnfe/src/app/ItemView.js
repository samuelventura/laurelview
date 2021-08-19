import React, { useReducer, useEffect } from 'react'

import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';

import ItemDisplay from "./ItemDisplay"

import socket from "../socket"
import env from "../environ"

function ItemView(props) {

  function reducer(state, { name, args, session }) {
    switch (name) {
      case "query": {
        const next = { ...state }
        next.query = args
        next.query.index = 0
        return next
      }
      case "close": {
        const next = { ...state }
        next.send = socket.send
        return next
      }
      case "send": {
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
    send: socket.send,
    query: {}
  }

  const [state, dispatch] = useReducer(reducer, initial)

  useEffect(() => {
    function handler({ name, args, session }) {
      dispatch({ name, args, session })
      switch (name) {
        case "send": {
          const name = "setup"
          const json = props.item.json
          const item = JSON.parse(json)
          const items = [item]
          args({ name, args: { items } })
          break
        }
        default: //linter complains
      }
    }
    return socket.createRt(handler, "/index")
  }, [props])

  function handleQuery(request) {
    const name = "query"
    const index = 0
    const response = ""
    const args = { index, request, response }
    state.send({ name, args })
  }

  function handleHide() {
    const action = "cancel"
    props.handler({ action })
  }

  return (
    <Modal show={props.show} onHide={handleHide} backdrop="static" centered>
      <Modal.Header closeButton>
        <Modal.Title>{props.item.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <ItemDisplay query={state.query} />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="success" onClick={() => handleQuery("read-value")}>Value</Button>
        <ButtonGroup>
          <Button variant="success" onClick={() => handleQuery("read-peak")}>Peak</Button>
          <Button variant="dark" onClick={() => handleQuery("reset-peak")}>Reset</Button>
        </ButtonGroup>
        <ButtonGroup>
          <Button variant="success" onClick={() => handleQuery("read-valley")}>Valley</Button>
          <Button variant="dark" onClick={() => handleQuery("reset-valley")}>Reset</Button>
        </ButtonGroup>
        <ButtonGroup>
          <Button variant="success" onClick={() => handleQuery("apply-tara")}>Tara</Button>
          <Button variant="dark" onClick={() => handleQuery("reset-tara")}>Reset</Button>
        </ButtonGroup>
        <Button variant="dark" onClick={() => handleQuery("reset-cold")}>Cold Reset</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ItemView
