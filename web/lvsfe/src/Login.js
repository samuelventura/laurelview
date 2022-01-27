import React, { useState } from 'react'

import api from "./api"
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Form from 'react-bootstrap/Form'
import Modal from 'react-bootstrap/Modal'
import Alert from 'react-bootstrap/Alert'
import Button from 'react-bootstrap/Button'

function Password(props) {

    const [currentPass, setCurrentPass] = React.useState("")

    //Response from yeico appliance
    const [responseString, setResponseString] = React.useState("")

    //Alerts
    const [isValid, setIsValid] = useState(false);
    const [isError, setIsError] = useState(false);

    //set default password in textbox so user does not have to copy-paste mac
    //this is because when reseting from physical button the locally cached 
    //password needs to be reset but there is only a global cache clear button
    //so user must copy mac to be able to login
    function buttonDefaultClick() {
        setCurrentPass(props.mac)
    }

    function buttonLoginClick() {
        var passEncode = Buffer.from(currentPass).toString('base64')
        api.getNetworkPing(function (res) {
            console.log(res)
            if (res.result === "ok") {
                setResponseString(`Login Success`)
                setIsValid(true)
                localStorage.setItem(props.mac, passEncode)
                setTimeout(() => {
                    setIsValid(false)
                    setCurrentPass("")
                }, 3000);
            }
            else {
                setResponseString(`Login Failed`)
                setIsError(true)
                setTimeout(() => {
                    setIsError(false)
                    setCurrentPass("")
                }, 3000);
            }
        }, props.device, "lvbox", currentPass)
    }

    return (
        <Modal
            {...props}
            size="lg"
            backdrop="static"
            aria-labelledby="contained-modal-title-vcenter"
            centered
        >
            <Modal.Header closeButton>
                <Modal.Title id="contained-modal-title-vcenter">
                    Login
                </Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form onSubmit={e => { e.preventDefault(); }}>
                    <Alert show={isValid} variant="success">
                        {responseString}
                    </Alert>
                    <Alert show={isError} variant="danger">
                        {responseString}
                    </Alert>

                    <Form.Group as={Row} className="mb-2">
                        <Form.Label align="right" column sm={3}>
                            Password
                        </Form.Label>
                        <Col sm={4}>
                            <Form.Control
                                type="password"
                                placeholder="Enter Current Pass"
                                value={currentPass}
                                onChange={e => setCurrentPass(e.target.value)}
                            /></Col>
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={buttonLoginClick} variant='dark'>Login</Button>
                <Button onClick={buttonDefaultClick} variant='dark'>Default</Button>
                <Button variant='dark' onClick={props.onHide}>Close</Button>
            </Modal.Footer>
        </Modal >
    );
}

export default Password;