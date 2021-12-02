import React, { useState, useContext } from 'react'
import Signer from './Signer'
import Api from './Api'

const initialSession = { id: "", email: "" }
const initialValue = { session: initialSession, setSession: () => { } }
const AuthContext = React.createContext(initialValue)

//based on https://kentcdodds.com/blog/how-to-use-react-context-effectively
function AuthProvider({ children }) {
  const [session, setSession] = useState(Api.fetchSession() || initialSession)
  // NOTE: you *might* need to memoize this value
  // Learn more in http://kcd.im/optimize-context
  const value = {
    session, setSession: (session) => {
      Api.saveSession(session)
      setSession(session)
    }
  }
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

function useAuth() {
  return useContext(AuthContext)
}

//https://gist.github.com/mjackson/d54b40a094277b7afdd6b81f51a0393f
function RequireAuth({ children }) {
  const auth = useAuth()
  return auth.session.id ? children : <Signer />
}

export { AuthProvider, RequireAuth, useAuth, initialSession }
