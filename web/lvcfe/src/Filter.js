import React, { useRef } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faSearch } from '@fortawesome/free-solid-svg-icons'

import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';

function Filter(props) {
  const filterInput = useRef(null);

  function onFilterKeyPress(e) {
    if (e.key === 'Enter') {
      onSearchClick()
    }
  }

  //Escape exits full screen on macos 
  //Escape not captured as keyPress
  function onFilterKeyUp(e) {
    if (e.key === 'Escape') {
      props.setFilter("")
    }
  }

  function onSearchClick() {
    props.setFilter(filterInput.current.value)
  }

  function onFilterChange(e) {
    props.setFilter(e.target.value)
  }

  function onClearClick() {
    props.setFilter("")
  }

  function onNewClick() {
    props.onNew()
  }

  return (<Form className="d-flex">
    <InputGroup>
      <InputGroup.Text><FontAwesomeIcon icon={faSearch} /></InputGroup.Text>
      <Form.Control value={props.filter} onChange={onFilterChange}
        onKeyPress={onFilterKeyPress} onKeyUp={onFilterKeyUp}
        placeholder="Filter..." type="text" ref={filterInput} />
      <Button onClick={onSearchClick} variant="outline-secondary" title="Apply Filter">Search</Button>
      <Button onClick={onClearClick} variant="outline-secondary" title="Clear Filter">Clear</Button>
      <Button onClick={onNewClick} variant="success" title="Create New">New...</Button>
    </InputGroup>
  </Form>)
}

export default Filter
