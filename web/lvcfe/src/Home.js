import React, { useState } from 'react'
import Table from 'react-bootstrap/Table';
import Filter from './Filter'
import Container from 'react-bootstrap/Container';

function Home() {
  const [filter, setFilter] = useState("")

  return (<div>

    <Container>
      <Filter filter={filter} setFilter={setFilter} />
      <Table striped bordered hover className="mt-3">
        <thead>
          <tr>
            <th>#</th>
            <th>Node Name</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Demo1</td>
          </tr>
          <tr>
            <td>2</td>
            <td>Demo2</td>
          </tr>
        </tbody>
      </Table>
    </Container>
  </div>)
}

export default Home
