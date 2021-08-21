import React, { useReducer, useEffect } from 'react'

import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCompress } from '@fortawesome/free-solid-svg-icons'
import { faExpand } from '@fortawesome/free-solid-svg-icons'

import ItemDisplay from "./ItemDisplay"

import useFullscreenStatus from "../fullscreen"
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
  
  const [isFullscreen, setIsFullscreen] = useFullscreenStatus(document.body)

  function toggleFullScreen() {
    setIsFullscreen(!isFullscreen)
  }

  function fullScreenButton() {
    //hide button if not supported
    if (document.fullscreenEnabled) {
      const icon = isFullscreen ? faCompress : faExpand
      return <Button variant="link" onClick={toggleFullScreen} 
        title="Toggle Full Screen"><FontAwesomeIcon icon={icon} /></Button>  
    }
  }

  function handleHide() {
    const action = "cancel"
    if (isFullscreen) {
      setIsFullscreen(false)
      //looks dark during transition
      setTimeout(()=>props.handler({ action }), 100)  
    } else {
      props.handler({ action })
    }
  }

  function query(index) {
    let q = state.queries[index]
    q = q || {}
    q.index = index
    return q
  }

  const displays = props.items.map((item, i) => {
    return <div key={item.id} className="front-panel">
        <h3>{item.name}</h3>
        <ItemDisplay query={query(i)} />
    </div>
  })

  return (
    <Modal show={props.show} onHide={handleHide} 
      backdrop="static" fullscreen centered>
      <Modal.Header closeButton bsPrefix="modal-header">
        <Modal.Title>{fullScreenButton()}Multi View</Modal.Title> 
      </Modal.Header>
      <Modal.Body>
        <div className="multi-grid">
        {displays}
        </div>
      </Modal.Body>
    </Modal>
  )
}

export default ItemMultiView
