import React, { useState, useEffect, useRef } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faArrowUp } from '@fortawesome/free-solid-svg-icons'
import { faArrowDown } from '@fortawesome/free-solid-svg-icons'
import { faSearch } from '@fortawesome/free-solid-svg-icons'

import Container from 'react-bootstrap/Container';
import Table from 'react-bootstrap/Table';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Navbar from 'react-bootstrap/Navbar';

import ItemEditor from "./ItemEditor"
import ItemDelete from "./ItemDelete"
import ItemControl from "./ItemControl"

import env from "../environ"

function ItemBrowser(props) {
  
  const [filter, setFilter] = useState("")
  const [sort, setSort] = useState("asc")

  const [showCreate, setShowCreate] = useState(false)
  const [showUpdate, setShowUpdate] = useState(false)
  const [showDelete, setShowDelete] = useState(false)
  const [showControl, setShowControl] = useState(false)
  const [itemSelected, setItemSelected] = useState({})

  const sortUpInput = useRef(null);
  const sortDownInput = useRef(null);
  const filterInput = useRef(null);

  useEffect(() => {
    //should not override validation focus
    if (filterInput.current != null) {
      filterInput.current.focus();
    }
    if (!props.state.online) {
      showDialog("none")
    }
  }, [props]);

  function onFilterKeyPress(e) {
    if(e.key === 'Enter') {
      onSearchClick()
    }
  }

  //esc exits full screen on macos 
  //use X icon to reset filter instead
  //Escape not captured as keyPress
  function onFilterKeyUp(e) {
    if(e.key === 'Escape') {
      setFilter("")
    }
  }

  function onSearchClick() {
    setFilter(filterInput.current.value)
  }

  function onFilterChange() {
    setFilter(filterInput.current.value)
  }

  function onClearClick() {
    setFilter("")
  }

  function handleSortChange(value) {
    //FIXME blur not working
    sortUpInput.current.blur()
    sortDownInput.current.blur()
    setSort(value)
  }

  function showDialog(action, item) {
    setShowCreate(false)
    setShowUpdate(false)
    setShowDelete(false)
    setShowControl(false)
    switch(action)
    {
      case "create":
        setShowCreate(true)
        break
      case "update":
        setItemSelected(item)
        setShowUpdate(true)
        break
      case "delete":
        setItemSelected(item)
        setShowDelete(true)
        break
      case "show":
        setItemSelected(item)
        setShowControl(true)
        break
      case "none":
        break
      default:
        env.log("Unknown action", action, item)
    }
  }

  function handleActions({action, args}) {
    switch(action) {
      case "cancel":
        showDialog("none")
        break;
      case "create": {
        const name = "create"
        props.dispatch({name, args})
        break;
      }
      case "update": {
        const name = "update"
        props.dispatch({name, args})
        break;
      }
      case "delete": {
        const name = "delete"
        props.dispatch({name, args})
        break;
      }
      default:
        env.log("Unknown action", action, args)
    }
  }

  function viewItems() {
    const f = filter.toLowerCase()
    const list = Object.values(props.state.items)
    const filtered = list.filter(item => 
      item.name.toLowerCase().includes(f))
    switch(sort) {
      case "asc":
        return filtered.sort((i1, i2) => 
          i1.name.localeCompare(i2.name))    
      case "desc":
        return filtered.sort((i1, i2) => 
          i2.name.localeCompare(i1.name))    
      default:
        return filtered 
    }     
  }

  const rows = viewItems().map(item => 
    <tr key={item.id}>
      <td>{item.name}</td>
      <td>
        <ButtonGroup>
        <Button onClick={() => showDialog("show", item)} 
          ref={sortUpInput} variant="outline-secondary" size="sm">Show</Button>        
        <Button onClick={() => showDialog("update", item)} 
          ref={sortUpInput} variant="outline-primary" size="sm">Edit</Button>        
        <Button onClick={() => showDialog("delete", item)} 
          ref={sortDownInput} variant="outline-danger" size="sm">Delete</Button>
        </ButtonGroup>
      </td>
    </tr>
  )

  function control(show) {
    if (show) {
      return <ItemControl show={showControl} item={itemSelected} handler={handleActions}/>
    }
  }

  return (
    <Container>

      <Navbar >
      <Navbar.Brand>Laurel View</Navbar.Brand>
      <Navbar.Toggle aria-controls="navbarScroll" />
      <Navbar.Collapse id="navbarScroll">

      <Form className="d-flex">
      <InputGroup>
      <InputGroup.Text><FontAwesomeIcon icon={faSearch} /></InputGroup.Text>
      <Form.Control value={filter} onChange={onFilterChange} 
            onKeyPress={onFilterKeyPress} onKeyUp={onFilterKeyUp} 
            placeholder="Filter..."  type="text" ref={filterInput}/>
      <Button onClick={onSearchClick} variant="outline-secondary">Search</Button>
      <Button onClick={onClearClick} variant="outline-secondary">Clear</Button>
      </InputGroup>
      </Form>
      </Navbar.Collapse>
      <Navbar.Collapse className="justify-content-end">
      <Button variant="success" onClick={() => showDialog("create")}>New...</Button>
      <ItemEditor show={showCreate} item={{}} handler={handleActions} action="create" title="Add New" button="Add New"/>
      <ItemEditor show={showUpdate} item={itemSelected} handler={handleActions} action="update" title="Update" button="Update"/>
      <ItemDelete show={showDelete} item={itemSelected} handler={handleActions} action="delete"/>
      { control(showControl) }
      </Navbar.Collapse>
      </Navbar>

      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Name {' '}
            <Button onClick={()=> handleSortChange("asc")} variant="link" size="sm">
              <FontAwesomeIcon icon={faArrowUp} /></Button>
            <Button onClick={()=> handleSortChange("desc")} variant="link" size="sm">
              <FontAwesomeIcon icon={faArrowDown} /></Button>
            </th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
        {rows}
        </tbody>
      </Table>
    </Container>
  )
}

export default ItemBrowser
