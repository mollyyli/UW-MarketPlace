import React from 'react';
import logo from './logo.svg';
import NavBar from "./NavBar/NavBar.js"
import { BrowserRouter as Router, Route, Link, NavLink, Switch } from "react-router-dom";

import './App.css';
import Listings from './Listings/Listings';
import Signin from './signin/Signin';
import ListingInfo from './ListingInfo/ListingInfo';
import SignUp from './SignUp/SignUp';
import AddListing from './AddListing/AddListing';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sid: ""
    }
  }
  handleStateChange = (newSid) => {
    // let sid = this.state.sid;
    this.setState({ sid: newSid });
    console.log(this.state.sid);
  }
  render() {
    console.log(this.state.sid);
    return (
      // <div className="App">
      <Router>
        <NavBar />
        <Route path="/listings" component={Listings} />
        <Route path="/sign-in" render={(props) => <Signin {...props} handleStateChange={this.handleStateChange} />} />

        <Route path="/listing/:listingID" component={ListingInfo} />
        <Route path="/sign-up" component={SignUp} />
        <Route path="/add" component={AddListing} />
      </Router>
      // </div>
    );
  }
}

export default App;
