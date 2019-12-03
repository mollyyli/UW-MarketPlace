import React from 'react';
import logo from './logo.svg';
import NavBar from "./NavBar/NavBar.js"
import { BrowserRouter as Router, Route, Link, NavLink, Switch } from "react-router-dom";

import './App.css';
import Listings from './Listings/Listings';
import Signin from './signin/Signin';
import ListingInfo from './ListingInfo/ListingInfo';
import AddListing from './AddListing/AddListing';

function App() {
  return (
    // <div className="App">
    <Router>
      <NavBar />
      <Route path="/listings" component={Listings} />
      <Route path="/sign-in" component={Signin} />

      <Route path="/listing/:listingID" component={ListingInfo} />
      <Route path="/add" component={AddListing} />
    </Router>
    // </div>
  );
}

export default App;
