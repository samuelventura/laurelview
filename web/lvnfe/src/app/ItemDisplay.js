import React from 'react'

import 'dseg/css/dseg.css';

function ItemDisplay(props) {

  const requestMap = {
    "read-value": "Value Reading",
    "read-peak": "Peak Reading",
    "reset-peak": "Resetting Peak",
    "read-valley": "Valley Reading",
    "reset-valley": "Resetting Valley",
    "apply-tara": "Applying Tara",
    "reset-tara": "Resetting Tara",
    "reset-cold": "Cold Resetting",
  }

  function requestText(request) {
    //falsy false, 0, "", null, undefined, NaN
    //truty "0"
    return requestMap[request] || "------"
  }

  function responseText(response) {
    return responseMap[response] || response || "------"
  }

  const responseMap = {
    "ok": "OK",
    "error": "Error",
  }

  return (
    <div className="display">
      <div className="request">{requestText(props.query.request)}</div>
      <div className="response">{responseText(props.query.response)}</div>
    </div>
  )
}

export default ItemDisplay
