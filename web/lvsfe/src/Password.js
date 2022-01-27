import React, { useState } from 'react'

import api from "./api"

import Modal from 'react-bootstrap/Modal'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Alert from 'react-bootstrap/Alert'

function Password(props) {

    const [newPass, setNewPass] = React.useState("")
    const [curPass, setCurPass] = React.useState(props.pass)

    React.useEffect(() => {
        setCurPass(props.pass);
    }, [props.pass])

    //Response from yeico appliance
    const [responseString, setResponseString] = React.useState("")

    //Alerts
    const [isValid, setIsValid] = useState(false);
    const [isError, setIsError] = useState(false);

    function buttonSetNewPassClick() {
        //console.log(props.pass, "->", curPass, "->", newPass)
        var passEncode = Buffer.from(newPass).toString('base64')
        api.setNewPass(function (res) {
            console.log(res)
            if (res.result === "ok") {
                setResponseString(`Set New Password Success`)
                setIsValid(true)
                setCurPass(newPass)
                setNewPass("")

                localStorage.setItem(props.mac, passEncode)
                setTimeout(() => {
                    setIsValid(false)
                }, 3000);
            }
            else {
                setResponseString(`Set New Password Fail`)
                setIsError(true)
                setTimeout(() => {
                    setIsError(false)
                    setNewPass("")
                }, 3000);
            }
        }, props.device, "lvbox", curPass, passEncode)
    }

    function buttonResetPassClick() {
        var passEncode = Buffer.from(props.mac).toString('base64')
        api.setDisablePass(function (res) {
            console.log(res)
            if (res.result === "ok") {
                setResponseString(`Reset Password Success`)
                setIsValid(true)
                localStorage.setItem(props.mac, passEncode)
                setTimeout(() => {
                    setIsValid(false)
                }, 3000);
            }
            else {
                setResponseString(`Reset Password Fail`)
                setIsError(true)
                setTimeout(() => {
                    setIsError(false)
                }, 3000);
            }
        }, props.device, "lvbox", curPass)
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
                    Password
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
                            New Password
                        </Form.Label>
                        <Col sm={8}>
                            <Form.Control
                                type="password"
                                placeholder="Enter New Password"
                                value={newPass}
                                onChange={e => setNewPass(e.target.value)}
                            />
                        </Col>
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={buttonSetNewPassClick} variant='dark'>Set New Password</Button>
                <Button onClick={buttonResetPassClick} variant='dark'>Reset Password</Button>
                <Button variant='dark' onClick={props.onHide}>Close</Button>
            </Modal.Footer>
        </Modal >
    );
}

export default Password;