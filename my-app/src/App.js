import React from 'react';
import logo from './logo.svg';
import NavBar from "./NavBar/NavBar.js"
import { BrowserRouter as Router, Route, Link, NavLink, Switch } from "react-router-dom";

import './App.css';
import Listings from './Listings/Listings';
import Signin from './signin/Signin';

function App() {
  return (
    // <div className="App">
    <Router>
      <NavBar />
      <Route path="/listings" component={Listings} />
      <Route path="/sign-in" component={Signin} />

    </Router>
    // </div>
  );
}

export default App;
