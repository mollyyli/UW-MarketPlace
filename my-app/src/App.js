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
import MyListings from './MyListings/MyListings'
import EditListings from './EditListings/EditListings'
import Splash from './Splash/Splash'

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      sid: localStorage.getItem('sid') ? localStorage.getItem('sid') : ""
    }
  }
  componentDidMount = async () => {
    let socket = new WebSocket(`wss://api.briando.me/v1/ws?auth=${this.state.sid}`);
    socket.onopen = () => {
      console.log("Websocket connection open")
    }
    socket.onclose = () => {
      console.log("Websocket connection closed")
    }
    socket.onmessage = (event) => {
      alert(event.data)
    }
    setTimeout(() =>  { alert("New listing") }, 60000);
  }

  handleStateChange = async (newSid) => {
    localStorage.setItem('sid', newSid);
    this.setState({ sid: newSid })
    console.log("handle change app state", this.state.sid);
  }
  render() {
    console.log("app state", this.state);
    return (
      <Router>
        <NavBar sid={this.state.sid} handleStateChange={this.handleStateChange} />
        <Route exact path="/" component={Splash} />
        <Route path="/listings" component={Listings} />
        <Route path="/sign-in" render={(props) => <Signin {...props} handleStateChange={this.handleStateChange} />} />

        <Route path="/listing/:listingID" component={ListingInfo} />
        <Route path="/sign-up" component={SignUp} />
        <Route path="/add" render={(props) => <AddListing {...props} sid={this.state.sid} />} />
        <Route path="/my-listings" render={(props) => <MyListings {...props} sid={this.state.sid} />} />
        <Route path="/edit/:listingID" render={(props) => <EditListings {...props} sid={this.state.sid} />}></Route>
      </Router>
      // </div>
    );
  }
}

export default App;
