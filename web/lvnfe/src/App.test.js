import { render, screen } from '@testing-library/react'
import App from './App'

test('renders filter... search box', () => {
  render(<App />)
  const searchBpx = screen.getByPlaceholderText(/Filter.../i)
  expect(searchBpx).toBeInTheDocument()
})
