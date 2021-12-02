import React, { useState, useContext } from 'react'
import { useLocation } from 'react-router-dom';

const initialAlert = { type: "", message: "" }
const initialValue = {
  current: initialAlert,
  clearAlert: () => { },
  errorAlert: () => { },
  warnAlert: () => { },
  successAlert: () => { }
}
const AlertContext = React.createContext(initialValue)

function AlertProvider({ children }) {
  const location = useLocation()
  const [alert, setAlert] = useState(initialAlert)
  // NOTE: you *might* need to memoize this value
  // Learn more in http://kcd.im/optimize-context
  const value = {
    current: alert,
    clearAlert: () => setAlert(initialAlert),
    errorAlert: (message) => setAlert({ type: "danger", message }),
    warnAlert: (message) => setAlert({ type: "warning", message }),
    successAlert: (message) => setAlert({ type: "success", message }),
  }
  React.useEffect(() => {
    setAlert(initialAlert)
  }, [location])
  React.useEffect(() => {
    if (alert.type) {
      console.log("alert", alert.type, alert.message)
    }
    if (alert.type === "success") {
      const timer = setTimeout(() => setAlert(initialAlert), 1000)
      return () => clearTimeout(timer)
    }
  }, [alert])
  return <AlertContext.Provider value={value}>{children}</AlertContext.Provider>
}

function useAlert() {
  return useContext(AlertContext)
}

export { AlertProvider, useAlert }
