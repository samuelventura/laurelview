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

import ItemMultiView from "./ItemMultiView"
import ItemEditor from "./ItemEditor"
import ItemDelete from "./ItemDelete"
import ItemView from "./ItemView"

import env from "../environ"

function ItemBrowser(props) {

  const [selected, setSelected] = useState({})
  const [filter, setFilter] = useState("")
  const [sort, setSort] = useState("asc")

  const [showCreate, setShowCreate] = useState(false)
  const [showUpdate, setShowUpdate] = useState(false)
  const [showDelete, setShowDelete] = useState(false)
  const [showView, setShowView] = useState(false)
  const [itemSelected, setItemSelected] = useState({})
  const [showMultiView, setShowMultiView] = useState(false)

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
    if (e.key === 'Enter') {
      onSearchClick()
    }
  }

  //esc exits full screen on macos 
  //use X icon to reset filter instead
  //Escape not captured as keyPress
  function onFilterKeyUp(e) {
    if (e.key === 'Escape') {
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
    sortUpInput.current.blur()
    sortDownInput.current.blur()
    setSort(value)
  }

  function showDialog(action, item) {
    setShowMultiView(false)
    setShowCreate(false)
    setShowUpdate(false)
    setShowDelete(false)
    setShowView(false)
    switch (action) {
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
      case "view":
        setItemSelected(item)
        setShowView(true)
        break
      case "multiview":
        setShowMultiView(true)
        break
      case "none":
        break
      default:
        env.log("Unknown action", action, item)
    }
  }

  function handleActions({ action, args }) {
    switch (action) {
      case "cancel":
        showDialog("none")
        break;
      case "create": {
        const name = "create"
        props.dispatch({ name, args })
        break;
      }
      case "update": {
        const name = "update"
        props.dispatch({ name, args })
        break;
      }
      case "delete": {
        const name = "delete"
        props.dispatch({ name, args })
        break;
      }
      default:
        env.log("Unknown action", action, args)
    }
  }

  function onSelectedChange(e, item) {
    const next = { ...selected }
    next[item.id] = e.target.checked
    setSelected(next)
  }

  function selectedItems() {
    const items = viewItems()
    return items.filter(item => selected[item.id])
  }

  function viewItems() {
    const f = filter.toLowerCase()
    const list = Object.values(props.state.items)
    const filtered = list.filter(item =>
      item.name.toLowerCase().includes(f))
    switch (sort) {
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
      <td>
        <Form.Check type="checkbox" label={item.name}
          value={item.selected} onChange={(e) => onSelectedChange(e, item)}></Form.Check>
      </td>
      <td>
        <ButtonGroup>
          <Button onClick={() => showDialog("view", item)}
            variant="link" size="sm">View</Button>
          <Button onClick={() => showDialog("update", item)}
            variant="link" size="sm">Edit</Button>
          <Button onClick={() => showDialog("delete", item)}
            variant="link" size="sm">Delete</Button>
        </ButtonGroup>
      </td>
    </tr>
  )

  function view(show) {
    if (show) {
      return <ItemView show={showView} item={itemSelected} handler={handleActions} />
    }
  }

  function multibutton() {
    const items = selectedItems()
    if (items.length > 0) {
      return <Button onClick={() => showDialog("multiview")}
        variant="link" size="sm">Multi View</Button>
    }
  }

  function multiview(show) {
    const items = selectedItems()
    if (show && items.length > 0) {
      return <ItemMultiView show={showMultiView} items={items} handler={handleActions} />
    }
  }

  return (
    <Container>

      <Navbar >
        <Navbar.Brand><img height="48px" src="banner.png" alt="Laurel View" /></Navbar.Brand>
        <Navbar.Collapse className="justify-content-end">
          <Button variant="success" onClick={() => showDialog("create")}>New...</Button>
        </Navbar.Collapse>
      </Navbar>

      <Navbar >
        <Navbar.Toggle aria-controls="navbarScroll" />
        <Navbar.Collapse id="navbarScroll">

          <Form className="d-flex">
            <InputGroup>
              <InputGroup.Text><FontAwesomeIcon icon={faSearch} /></InputGroup.Text>
              <Form.Control value={filter} onChange={onFilterChange}
                onKeyPress={onFilterKeyPress} onKeyUp={onFilterKeyUp}
                placeholder="Filter..." type="text" ref={filterInput} />
              <Button onClick={onSearchClick} variant="outline-secondary">Search</Button>
              <Button onClick={onClearClick} variant="outline-secondary">Clear</Button>
            </InputGroup>
          </Form>
        </Navbar.Collapse>
        <Navbar.Collapse className="justify-content-end">
          <ItemEditor show={showCreate} item={{}} handler={handleActions} action="create" title="Add New" button="Add New" />
          <ItemEditor show={showUpdate} item={itemSelected} handler={handleActions} action="update" title="Update" button="Update" />
          <ItemDelete show={showDelete} item={itemSelected} handler={handleActions} action="delete" />
          {view(showView)}
          {multiview(showMultiView)}
        </Navbar.Collapse>
      </Navbar>

      <Table striped bordered hover>
        <thead>
          <tr>
            <th>Name &nbsp;
              <Button ref={sortUpInput} onClick={() => handleSortChange("asc")} variant="link" size="sm">
                <FontAwesomeIcon icon={faArrowUp} /></Button>
              <Button ref={sortDownInput} onClick={() => handleSortChange("desc")} variant="link" size="sm">
                <FontAwesomeIcon icon={faArrowDown} /></Button>
            </th>
            <th>
              Actions {multibutton()}
            </th>
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
