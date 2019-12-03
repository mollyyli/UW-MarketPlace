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

class NavBar extends Component {
  render() {
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
            <NavItem>
              <NavLink href="/sign-in">Sign In</NavLink>
            </NavItem>
            <NavItem>
              <NavLink href="/sign-up">Sign Up</NavLink>
            </NavItem>
          </Nav>
        </Navbar>
        </header>
      </div>
    );
  }
}

export default NavBar;
