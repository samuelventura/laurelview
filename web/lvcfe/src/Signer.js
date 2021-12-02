import React from 'react'
import Tabs from 'react-bootstrap/Tabs'
import Tab from 'react-bootstrap/Tab'
import SignIn from './SignIn'
import SignUp from './SignUp'
import Recover from './Recover'

function Signer() {
  return (
  <Tabs defaultActiveKey="signin" className="mb-3 mx-autO" md="auto">
    <Tab eventKey="signin" title="Sign-in">
      <SignIn/>
    </Tab>
    <Tab eventKey="signup" title="Sign-up">
      <SignUp/>
    </Tab>
    <Tab eventKey="recover" title="Recover">
      <Recover/>
    </Tab>
  </Tabs>
  )
}

export default Signer
