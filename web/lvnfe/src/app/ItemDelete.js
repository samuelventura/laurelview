import React from 'react'

import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';

function ItemDelete(props) {

  function handleHide() {
    const action = "cancel"
    props.handler({ action })
  }

  function handleAction() {
    const action = "delete"
    const id = props.item.id
    const args = { id }
    props.handler({ action, args })
    handleHide()
  }

  return (
    <Modal show={props.show} onHide={handleHide} centered>
      <Modal.Header closeButton>
        <Modal.Title>Caution</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <p>Delete item <b>{props.item.name}</b>?</p>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleHide}>Close</Button>
        <Button variant="danger" onClick={handleAction}>Delete</Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ItemDelete
