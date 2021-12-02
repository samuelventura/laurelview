import { sha256 } from 'js-sha256';

const logit = function(data) {
    console.log("fetch log:", data)
}

function fetchIt(url, params) {
    const logit = function(data) {
        console.log("fetch log:", url, params, data)
    }
    return new Promise((resolve, reject) => {
        console.log("fetching", url, params)
        //stringify and upper case all errors
        const rejectIt = msg => reject(`${msg}`.toUpperCase())
        fetch(url, params)
        .then(res => {
            switch (res.status) {
                case 200:
                    res.json().then(resolve).catch(logit)
                    break
                case 400:
                    res.text().then(rejectIt).catch(logit)
                    break
                default:
                    res.text().then(logit).catch(logit)
                    break
                }
        })
        .catch(logit)
    })
}

function signup(email) {
    const query = new URLSearchParams();
    query.set("aid", email)
    const url = `api/signup?${query.toString()}`
    const method = "POST"
    return fetchIt(url, {method})
}

function signin(email, password) {
    const query = new URLSearchParams();
    query.set("aid", email)
    const url = `api/signin?${query.toString()}`
    const method = "POST"
    const body = sha256(password);
    return fetchIt(url, {method, body})
}

function signout(sid) {
    localStorage.removeItem("lv.session")
    const query = new URLSearchParams();
    query.set("sid", sid)
    const url = `api/signout?${query.toString()}`
    const method = "GET"
    fetchIt(url, {method}).then(logit).catch(logit)
}

function recover(email) {
    const query = new URLSearchParams();
    query.set("aid", email)
    const url = `api/recover?${query.toString()}`
    const method = "POST"
    return fetchIt(url, {method})
}

function fetchSession() {
    const session = localStorage.getItem("lv.session")
    return session ? JSON.parse(session) : null
}

function saveSession(session) {
    const json = JSON.stringify(session)
    return localStorage.setItem("lv.session", json)
}

const exports = { 
    fetchSession,
    saveSession,
    signup,
    signin,
    signout,
    recover,
}

export default exports
