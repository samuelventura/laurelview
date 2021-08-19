import React, { useReducer, useEffect } from 'react'

import 'dseg/css/dseg.css';
import env from "../environ"

function ItemDisplay(props) {

  function reducer(state, { name, args }) {
    switch (name) {
      case "millis": {
        const next = { ...state }
        next.millis = args
        return next
      }
      case "query": {
        const next = { ...state }
        next.latency = next.millis
        next.start = new Date()
        next.query = args
        next.millis = 0
        return next
      }
      default:
        env.log("Unknown mutation", name, args)
        return state
    }
  }

  const initial = {
    start: new Date(),
    latency: 0,
    millis: 0,
    query: {}
  }

  const [state, dispatch] = useReducer(reducer, initial)

  useEffect(() => {
    const tid = setInterval(()=>{
      const now = new Date().getTime()
      const elapsed = now - state.start.getTime()
      dispatch({name:"millis", args: elapsed})
    }, 10)
    return () => clearInterval(tid)
  }, [state.start])

  useEffect(() => {
    const q = props.query || {}
    dispatch({name:"query", args: q})
    if (q.error) {
      console.error(q.index, q.error, q)
    }
  }, [props.query])

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
    return `${state.millis} @ ${state.latency} ms`
  }

  return (
    <div className="display">
      <div className="top">
        <div className="timing">{timing()}</div>
        <div className="request">{requestText(state.query.request)}</div>
      </div>
      <div className="response">{responseText(state.query.response)}</div>
      <div className="error">{state.query.error}</div>
    </div>
  )
}

export default ItemDisplay
