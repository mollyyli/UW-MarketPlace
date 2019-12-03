import React, { Component } from 'react';
import { Container, Alert } from 'reactstrap';

import { Link } from 'react-router-dom';
import "./AddListing.css";
import '../App.css';
import AddNameForm from './AddListingForm';

class AddListing extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: "",
      description: "",
      condition: "",
      location: "",
      price: "",
      contact: "",
      validated: true
    }
  }
  componentDidMount = () => {
    if (!this.props.sid) {
      this.props.history.push("/sign-in")
    }
  }
  handleInputChange = ({ target: { name, value } }) => {
    this.setState({ [name]: value });
  }
  handleSubmit = async event => {
    event.preventDefault();
    if (this.state.title && this.state.description && this.state.condition && this.state.price && this.state.contact) {
      this.setState({ validated: true });
      const response = await fetch('https://api.briando.me/v1/listings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': this.props.sid
        },
        body: JSON.stringify({
          "title": this.state.title,
          "description": this.state.description,
          "condition": this.state.condition,
          "price": this.state.price,
          "contact": this.state.contact,
          "location": this.state.location
        })
      })
      const listing = await (response.json());
      console.log(listing);
      this.props.history.push(`/listing/${listing._id}`)
      console.log(this.state);
    } else {
      this.setState({ validated: false });
    }
  }
  render() {
    return (
      <div className="ListingInfo">
        <Container>
          <h1>Add a listing</h1>
          <AddNameForm
            {...this.state}
            onHandleInputChange={this.handleInputChange}
            onHandleSubmit={this.handleSubmit}
          />
          {!this.state.validated &&
            <Alert className="validate-alert" color="danger">
              Please fill out all fields
            </Alert>
          }
        </Container>
      </div>
    );
  }
}

export default AddListing;
