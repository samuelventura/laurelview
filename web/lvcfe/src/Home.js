import React from 'react'
import Table from 'react-bootstrap/Table'

function Home() {
  return (<div>
    <Table striped bordered hover>
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
  </div>)
}

export default Home
