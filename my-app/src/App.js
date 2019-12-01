import React from 'react';
import logo from './logo.svg';
import NavBar from "./NavBar/NavBar.js"
import { BrowserRouter as Router, Route, Link, NavLink, Switch } from "react-router-dom";

import './App.css';
import Listings from './Listings/Listings';

function App() {
  return (
    // <div className="App">
    <Router>
      <NavBar />
      <Route path="/listings" component={Listings} />
    </Router>
    // </div>
  );
}

export default App;
