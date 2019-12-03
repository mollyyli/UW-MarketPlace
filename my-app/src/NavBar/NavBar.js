import React, { Component, useState } from 'react';
import {
  Collapse,
  Navbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink,
  UncontrolledDropdown,
  DropdownToggle,
  DropdownMenu,
  DropdownItem
} from 'reactstrap';
import '../App.css';
import './NavBar.css';

class NavBar extends Component {
  signOut = async () => {
    await fetch('https://api.briando.me/v1/sessions/mine', {
      method: 'DELETE',
      headers: {
        'Authorization': this.props.sid
      }
    })
    this.props.handleStateChange("");
  }
  render() {
    console.log("sid", this.props.sid);
    return (
      <div className="NavBar">
        <header>
          <Navbar color="light" light expand="md">
            <NavbarBrand href="/">UW Marketplace</NavbarBrand>
            <Nav className="mr-auto" navbar>
              <NavItem>
                <NavLink href="/listings">Listings</NavLink>
              </NavItem>
              <NavItem>
                <NavLink href="/add">Add Listing</NavLink>
              </NavItem>
              {!this.props.sid ?
                <span className="nav-span">
                  <NavItem>
                    <NavLink href="/sign-in">Sign In</NavLink>
                  </NavItem>
                  <NavItem>
                    <NavLink href="/sign-up">Sign Up</NavLink>
                  </NavItem>
                </span> :
                <span className="nav-span">
                  <NavItem>
                    <NavLink href="/my-listings">My Listings</NavLink>
                  </NavItem>
                  <NavItem>
                    <NavLink onClick={this.signOut}>Sign Out</NavLink>
                  </NavItem>
                </span>
              }
            </Nav>
          </Navbar>
        </header>
      </div>
    );
  }
}

export default NavBar;
