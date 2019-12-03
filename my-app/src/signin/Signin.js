import React, { Component } from 'react';
import { Button, FormGroup, FormControl } from "react-bootstrap";

import "./Signin.css";
import '../App.css';
// import { sign } from 'crypto';

class Signin extends Component {
  constructor(props) {
    super(props);
    this.state = {
      // listings: [],
      // loading: false,
      email: '',
      password: '',
    }
  }
  // componentDidMount = async () => {
  //   //await this.validateForm();
  //   // await this.handleSubmit();
  //  await this.signin();
  // }

  validateForm = async () =>  {
    return this.state.email.length > 0 && this.state.password.length > 0;
  }

  handleSubmit = async (event) => {
    let body = {
      email: this.state.email,
      password: this.state.password,
    }
    event.preventDefault();
    const response = await fetch('https://api.briando.me/v1/sessions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',

      },
      body: JSON.stringify(body),
    });
    const signin = await response.json();
    this.props.handleStateChange(response.headers.get("Authorization"));
    this.props.history.push("/listings"); 
  }

  handleEmailChange = (event) => {
    this.setState({
      email: event.target.value,
    });
  }

  handlePasswordChange = (event) => {
    this.setState({
      password: event.target.value,
    });
  }

  signin = async () => {
    // this.setState({ loading: true });
    
    // this.setState({ listings: listings, loading: false });
    
  }

  render() {
    return (
      <div className="signin">
        <form className="email" onSubmit={this.handleSubmit}>
          <FormGroup controlId="email">
            <label>Email</label>
            <FormControl
              autoFocus
              type="email"
              value={this.state.email}
              // onChange={e => setEmail(e.target.value)}
              onChange={this.handleEmailChange}

            />
          </FormGroup>

          <FormGroup className="password" controlId="password">
            <label>Password</label>
            <FormControl
              value={this.state.password}
              // onChange={e => setPassword(e.target.value)}
              onChange={this.handlePasswordChange}
              type="password"
            />
          </FormGroup>

      
          <Button className="button" block disabled={!this.validateForm()} type="submit">
            Sign in
          </Button>
        </form>
      </div>
    );
  }
}

export default Signin;
