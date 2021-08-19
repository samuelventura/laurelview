import React, { useState, useEffect } from 'react'

import 'dseg/css/dseg.css';

function ItemDisplay(props) {

  const [query, setQuery] = useState({});
  const [start, setStart] = useState(new Date());
  const [millis, setMillis] = useState(0);
  const [latency, setLatency] = useState(0);

  useEffect(() => {
    const tid = setInterval(()=>{
      setMillis(elapsed())
    }, 1)
    return () => clearInterval(tid)
  }, [start])

  useEffect(() => {
    setLatency(elapsed())
    setQuery(props.query)
    setStart(new Date())
    setMillis(0)
  }, [props.query])

  function elapsed() {
    const now = new Date().getTime()
    return now - start.getTime()
  }

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

  const responseMap = {
    "ok": "OK",
    "error": "Error",
  }

  function requestText(request) {
    //falsy false, 0, "", null, undefined, NaN
    //truty "0"
    return requestMap[request] || "------"
  }

  function responseText(response) {
    return responseMap[response] || response || "------"
  }

  function timing() {
    let str = `${latency} ms`
    if (millis > 1000) {
      str = `${millis} ! ` + str
    }
    return str
  }

  return (
    <div className="display">
      <div className="top">
        <div className="timing">{timing()}</div>
        <div className="request">{requestText(query.request)}</div>
      </div>
      <div className="response">{responseText(query.response)}</div>
      <div className="error">{query.error}</div>
    </div>
  )
}

export default ItemDisplay
