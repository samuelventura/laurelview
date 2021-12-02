import React, { useState } from 'react'
import Form from 'react-bootstrap/Form'
import InputGroup from 'react-bootstrap/InputGroup'
import Button from 'react-bootstrap/Button'
import { useAlert } from './Alert'
import Api from './Api'

function Recover() {
  const alert = useAlert()
  const [email, setEmail] = useState("")
  const [validated, setValidated] = useState(false);
  function handleSubmit(e) {
    e.preventDefault()
    e.stopPropagation()
    setValidated(true)
    const form = e.currentTarget;
    if (!form.checkValidity()) return
    alert.warnAlert("Recovering...")
    Api.recover(email)
      .then(msg => alert.successAlert(msg))
      .catch(error => alert.errorAlert(error))
  }
  return (
    <React.Fragment>
      <Form noValidate validated={validated} onSubmit={handleSubmit}>
        <Form.Group className="mb-3" controlId="formRecoverEmail">
          <Form.Label>Email</Form.Label>
          <InputGroup hasValidation>
            <Form.Control type="email" placeholder="Email" required
              value={email} onChange={e => setEmail(e.target.value)} />
            <Form.Control.Feedback type="invalid">Provide a valid email.</Form.Control.Feedback>
          </InputGroup>
        </Form.Group>
        <Button variant="primary" type="submit" className="float-end">Recover</Button>
      </Form>
    </React.Fragment>
  )
}

export default Recover
