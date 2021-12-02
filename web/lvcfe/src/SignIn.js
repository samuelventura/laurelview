import React, { useState } from 'react'
import Form from 'react-bootstrap/Form'
import InputGroup from 'react-bootstrap/InputGroup'
import Button from 'react-bootstrap/Button'
import { useAuth } from './Auth'
import { useAlert } from './Alert'
import Api from './Api'

function SignIn() {
  const auth = useAuth()
  const alert = useAlert()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [validated, setValidated] = useState(false);
  function handleSubmit(e) {
    e.preventDefault()
    e.stopPropagation()
    setValidated(true)
    const form = e.currentTarget;
    if (!form.checkValidity()) return
    alert.warnAlert("Signing in...")
    Api.signin(email, password)
      .then(s => auth.setSession({ id: s.Sid, email: s.Aid }))
      .then(() => alert.successAlert("Success signing in"))
      .catch(error => alert.errorAlert(error))
  }
  return (
    <React.Fragment>
      <Form noValidate validated={validated} onSubmit={handleSubmit}>
        <Form.Group className="mb-3" controlId="formSignInEmail">
          <Form.Label>Email</Form.Label>
          <InputGroup hasValidation>
            <Form.Control type="email" placeholder="Email" required
              value={email} onChange={e => setEmail(e.target.value)} />
            <Form.Control.Feedback type="invalid">Provide a valid email.</Form.Control.Feedback>
          </InputGroup>
        </Form.Group>
        <Form.Group className="mb-3" controlId="formSignInPassword">
          <Form.Label>Password</Form.Label>
          <InputGroup hasValidation>
            <Form.Control type="password" placeholder="Password" required
              value={password} onChange={e => setPassword(e.target.value)} />
            <Form.Control.Feedback type="invalid">Provide a non empty password.</Form.Control.Feedback>
          </InputGroup>
        </Form.Group>
        <Button variant="primary" type="submit" className="float-end">Sign-in</Button>
      </Form>
    </React.Fragment>
  )
}

export default SignIn
