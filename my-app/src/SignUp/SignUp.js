import React, { Component } from 'react';
import {
  Card, Container, Row, Button, Spinner, Form, FormGroup, Label, Input, Alert
} from 'reactstrap';
import { Link } from 'react-router-dom';
import '../App.css';
import './SignUp.css';

class SignUp extends Component {
  constructor(props) {
    super(props);
    this.state = {
      success: false,
      signUpAttempt: false,
      loading: false
    }
  }

  formSubmit = async (data) => {
    data.preventDefault()
    this.setState({loading: true});
    const response = await fetch(`https://api.briando.me/v1/users`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "Email": `${data.target.email.value}`,
            "Password": `${data.target.password.value}`,
            "PasswordConf": `${data.target.passwordconf.value}`,
            "UserName": `${data.target.username.value}`,
            "FirstName": `${data.target.firstname.value}`,
            "LastName": `${data.target.lastname.value}`
        })
    });
    const listing = await response;
    this.setState({loading: false});
    if (listing.status == 201) {
        this.setState({success: true})
    } else {
        this.setState({signUpAttempt: true})
    }
    }
  render() {
    let content, signup;
    if (this.state.success) {
        content = 
            <Alert color="success">
                Successfully Signed Up!
            </Alert>
    } else {
        content = 
        <Form onSubmit={this.formSubmit}>
            <FormGroup>
            <Label for="exampleEmail">Email</Label>
            <Input type="email" name="email" id="exampleEmail" placeholder="email@email.com" />
            </FormGroup>
            <FormGroup>
            <Label for="password">Password</Label>
            <Input type="password" name="password" id="password" placeholder="password" />
            </FormGroup>
            <FormGroup>
            <Label for="passwordconf">Password Confirmation</Label>
            <Input type="password" name="passwordconf" id="passwordconf" placeholder="password confirmation" />
            </FormGroup>
            <FormGroup>
            <Label for="username">Username</Label>
            <Input name="username" id="username" placeholder="user name" />
            </FormGroup>
            <FormGroup>
            <Label for="firstname">First Name</Label>
            <Input name="firstname" id="firstname" placeholder="first name" />
            </FormGroup>
            <FormGroup>
            <Label for="lastname">Last Name</Label>
            <Input name="lastname" id="lastname" placeholder="last name" />
            </FormGroup>
            <Button>{this.state.loading ? <Spinner size="sm" color="light" />: "Submit"}</Button>
        </Form>;
        if (this.state.signUpAttempt) {
            signup = 
                <Alert color="danger">
                    Failed to Sign Up
                </Alert>
        }
        
    }
    return (
        <Container className="sign-up-container">
            <h1>Sign Up</h1>
            {signup}
            {content}
        </Container>
    );
  }
}

export default SignUp;
