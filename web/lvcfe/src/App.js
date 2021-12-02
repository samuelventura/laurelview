import React from 'react'
import { Outlet, useNavigate } from "react-router-dom"
import Navbar from 'react-bootstrap/Navbar'
import Container from 'react-bootstrap/Container'
import Button from 'react-bootstrap/Button'
import Alert from 'react-bootstrap/Alert'
import {useAuth, initialSession} from './Auth'
import {useAlert} from './Alert'
import Api from './Api'
import './App.css'

function Username() {
  const auth = useAuth()
  const sid = auth.session.id
  const handleSigoutClick = function() {
    auth.setSession(initialSession)
    Api.signout(sid)
  } 
  if (sid) return (
    <Button variant="link" onClick={handleSigoutClick} title="Sign-out">
      {auth.session.email}
    </Button> 
  )
  return null
}

function showAlert(alert) {
  const current = alert.current
  if (current.type) {
    return (<Alert variant={current.type} dismissible
      onClose={() => alert.clearAlert()}>
      {current.message}
      </Alert>)
  }
}

//App component acts like app layout
//Navbar.Brand made pointer cursor with class btn
function App() {
  const alert = useAlert()
  const navigate = useNavigate()
  const handleBrandClick = () => navigate("/", { replace: true })
  return (
    <Container>
      <Navbar className="mb-3">
        <Container>
          <Navbar.Brand  className="btn" onClick={handleBrandClick} title="Home">
            <img height="48px" src="banner.png" alt="Home"/>
          </Navbar.Brand>
          <Navbar.Toggle />
          <Navbar.Collapse className="justify-content-end">
            <Navbar.Text><Username/></Navbar.Text>
          </Navbar.Collapse>
        </Container>        
      </Navbar>
      <Container>
        {showAlert(alert)}
        <Outlet />
      </Container> 
    </Container> 
  )
}

export default App
