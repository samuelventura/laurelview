import React, { useState, useEffect, useRef } from 'react'

import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';

//FIXME enter triggers navigation and validation
function ItemEditor(props) {
  
  const nameInput = useRef(null);
  const hostInput = useRef(null);
  const portInput = useRef(null);
  const slaveInput = useRef(null);

  useEffect(() => {
    //should not override validation focus
    if (nameInput.current != null) {
      nameInput.current.focus();
    }
    const json = props.item.json || "{}"
    const args = JSON.parse(json)
    setName(props.item.name || "")
    setHost(args.host || "")
    setPort(args.port || "0")
    setSlave(args.slave || "1")

    //red flicker on disconnection while editing
    setValidated(false);
  }, [props]);

  const [name, setName] = useState("")
  const [host, setHost] = useState("")
  const [port, setPort] = useState("")
  const [slave, setSlave] = useState("")

  const [validated, setValidated] = useState(false);

  function handleHide() {
    setValidated(false);
    const action = "cancel"
    props.handler({action})
  }

  function handleAction() {
    var errors = 0
    if (!nameInput.current.checkValidity()) {
      nameInput.current.focus()
      errors++
    } else if (!hostInput.current.checkValidity()) {
      hostInput.current.focus()
      errors++
    } else if (!portInput.current.checkValidity()) {
      portInput.current.focus()
      errors++
    } else if (!slaveInput.current.checkValidity()) {
      slaveInput.current.focus()
      errors++
    }
    if (errors > 0) {
      setValidated(true);
      return
    }
    const action = props.action
    const data = {host, port, slave}
    const json = JSON.stringify(data)
    const id = props.item.id || 0
    const args = {id, name, json}
    props.handler({action, args})
    handleHide()
  }

  function onNameChange(e) {
    setName(e.target.value)
  }

  function onHostChange(e) {
    setHost(e.target.value)
  }  

  function onPortChange(e) {
    setPort(e.target.value)
  }  

  function onSlaveChange(e) {
    setSlave(e.target.value)
  }

  return (
    <Modal show={props.show} onHide={handleHide} centered>
      <Modal.Header closeButton>
        <Modal.Title>{props.title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
      <Form validated={validated}>
        <Form.Group className="mb-3" controlId="itemName">
          <Form.Label>Name</Form.Label>
          <Form.Control value={name} onChange={onNameChange} 
            type="text" placeholder="Name" pattern="\S(.*\S)?" 
            required ref={nameInput}/>
          <Form.Control.Feedback type="invalid">
            Name cannot be blank nor have white spaces on the edges
            </Form.Control.Feedback>            
        </Form.Group>

        <Form.Group className="mb-3" controlId="itemHost">
          <Form.Label>Host</Form.Label>
          <Form.Control value={host} onChange={onHostChange} 
            type="text" placeholder="Host" pattern="\S+" 
            required ref={hostInput}/>
          <Form.Control.Feedback type="invalid">
            Host cannot be blank nor have white spaces on the edges
            </Form.Control.Feedback>
        </Form.Group>

        <Form.Group className="mb-3" controlId="itemPort">
          <Form.Label>Port</Form.Label>
          <Form.Control value={port} onChange={onPortChange} 
            type="number" placeholder="Port" required 
            min="1" max="65535" ref={portInput}/>
          <Form.Control.Feedback type="invalid">
            Port must be an integer between 1 and 65535
            </Form.Control.Feedback>            
        </Form.Group>

        <Form.Group className="mb-3" controlId="itemSlave">
          <Form.Label>Slave</Form.Label>
          <Form.Control value={slave} onChange={onSlaveChange} 
            type="number" placeholder="Slave" required 
            min="1" max="31" ref={slaveInput}/>
          <Form.Control.Feedback type="invalid">
            Slave must be an integer between 1 and 31
            </Form.Control.Feedback>            
        </Form.Group>
      </Form>
      </Modal.Body>

      <Modal.Footer>
        <Button variant="secondary" onClick={handleHide}>Close</Button>
        <Button variant="primary" onClick={handleAction}>{props.button}</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ItemEditor
