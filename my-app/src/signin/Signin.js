import React, { Component } from 'react';
import { Button, Input, FormGroup, Alert, Container, Spinner } from "reactstrap";

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
      error: false,
      loading: false
    }
  }
  // componentDidMount = async () => {
  //   //await this.validateForm();
  //   // await this.handleSubmit();
  //  await this.signin();
  // }

  validateForm = async () => {
    return this.state.email.length > 0 && this.state.password.length > 0;
  }

  handleSubmit = async (event) => {
    this.setState({ loading: true });
    // console.log(this.state.email, this.state.password);
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
    response.status !== 201 && this.setState({ error: true, loading: false })
    const signin = await response.json();
    this.props.handleStateChange(response.headers.get("Authorization"));
    this.setState({ loading: false });
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
        <Container>
          <h1>Sign In</h1>
          <form className="email" onSubmit={this.handleSubmit}>
            <FormGroup controlId="email">
              <label>Email</label>
              <Input
                autoFocus
                type="email"
                value={this.state.email}
                // onChange={e => setEmail(e.target.value)}
                onChange={this.handleEmailChange}

              />
            </FormGroup>

            <FormGroup className="password" controlId="password">
              <label>Password</label>
              <Input
                value={this.state.password}
                // onChange={e => setPassword(e.target.value)}
                onChange={this.handlePasswordChange}
                type="password"
              />
            </FormGroup>


            <Button className="button" block disabled={!this.validateForm()} type="submit">
              {!this.state.loading ? "Sign in" : <Spinner size="sm" color="light" />}
            </Button>
            {this.state.error && <Alert className="email-alert" color="danger">Email doesn't exist or password is incorrect.</Alert>}
          </form>
        </Container>
      </div>
    );
  }
}

export default Signin;
