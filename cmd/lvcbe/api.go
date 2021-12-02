package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samuelventura/go-tree"
)

type apiDso struct {
	dao  *daoDso
	sid  *idDso
	pid  *idDso
	stom int64
}

func newApi(node tree.Node) *apiDso {
	dao := node.GetValue("dao").(*daoDso)
	stom := node.GetValue("stom").(int64)
	sid := newId("sid")
	pid := newId("pid")
	return &apiDso{dao, sid, pid, stom}
}

//FIXME add delayer to prevent attacks
func (dso *apiDso) post_signup(id string) (*AccountDro, error) {
	password := pwdit(dso.pid.next())
	dro := &AccountDro{}
	dro.Aid = id
	dro.Created = time.Now()
	dro.Password = hashit(dso.pid.next())
	dro.Recover = hashit(password)
	dro.Enabled = true
	err := dso.dao.create(dro)
	if err != nil {
		err = fmt.Errorf("account already exists")
		return nil, err
	}
	//FIXME remove development log
	log.Println("signup", id, password)
	//FIXME send email with recover password
	return dro, err
}

//FIXME add delayer to prevent attacks
func (dso *apiDso) post_signin(aid, hash string) (*SessionDro, error) {
	adro, err := dso.dao.getAccount(aid)
	if err != nil {
		err = fmt.Errorf("account not found")
		return nil, err
	}
	if !adro.Enabled {
		err = fmt.Errorf("account is disabled")
		return nil, err
	}
	//log.Println("signin", aid, hash, adro.Password, adro.Recover)
	//FIXME never leave neither password nor recover empty
	if hash != adro.Password && hash != adro.Recover {
		err = fmt.Errorf("invalid credentials")
		return nil, err
	}
	//recover password usable only once
	if hash == adro.Recover {
		adro.Password = adro.Recover
		adro.Recover = hashit(dso.pid.next())
		err = dso.dao.update(adro)
		if err != nil {
			return nil, err
		}
	}
	stom := time.Duration(dso.stom)
	sdro := &SessionDro{}
	sdro.Sid = hashit(dso.sid.next())
	sdro.Aid = aid
	sdro.Created = time.Now()
	sdro.Expires = time.Now().Add(stom * time.Minute)
	sdro.Enabled = true
	err = dso.dao.create(sdro)
	return sdro, err
}

func (dso *apiDso) get_signout(sid string) (*SessionDro, error) {
	sdro, err := dso.dao.getSession(sid)
	if err != nil {
		err = fmt.Errorf("session not found")
		return nil, err
	}
	err = dso.dao.delete(sdro)
	if err != nil {
		err = fmt.Errorf("session delete error")
		return nil, err
	}
	return sdro, nil
}

//FIXME add delayer to prevent attacks
func (dso *apiDso) post_recover(aid, hash string) (*AccountDro, error) {
	dro, err := dso.dao.getAccount(aid)
	if err != nil {
		err = fmt.Errorf("account not found")
		return nil, err
	}
	if !dro.Enabled {
		err = fmt.Errorf("account is disabled")
		return nil, err
	}
	password := pwdit(dso.pid.next())
	dro.Recover = hashit(password)
	err = dso.dao.update(dro)
	if err != nil {
		return nil, err
	}
	//FIXME remove development log
	log.Println("recover", aid, password)
	return dro, err
}
