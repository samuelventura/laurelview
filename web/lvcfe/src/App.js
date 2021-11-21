import logo from './logo.png';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          LaurelView Live Demos
        </p>
        <a
          className="App-link"
          href="/proxy/demo1/"
        >
          Live Demo1
        </a>
        <a
          className="App-link"
          href="/proxy/demo2/"
        >
          Live Demo2
        </a>
      </header>
    </div>
  );
}

export default App;
