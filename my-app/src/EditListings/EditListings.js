import React, { Component } from 'react';
import { Container, Alert } from 'reactstrap';

import { Link } from 'react-router-dom';
import '../App.css';
import EditListingForm from './EditListingForm';

class EditListings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: "",
      description: "",
      condition: "",
      price: "",
      contact: "",
      location: "",
      validated: true
    }
  }
  handleInputChange = ({ target: { name, value } }) => {
    this.setState({ [name]: value });
  }
  handleSubmit = async event => {
    event.preventDefault();
    if (this.state.title && this.state.description && this.state.condition && this.state.price && this.state.contact && this.state.location) {
      this.setState({ validated: true });
      const listingID = this.props.match.params.listingID;
      await fetch(`https://api.briando.me/v1/listings/${listingID}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          "Authorization": this.props.sid
        },
        body: JSON.stringify({
          "title": this.state.title,
          "description": this.state.description,
          "condition": this.state.condition,
          "price": this.state.price,
          "contact": this.state.contact,
          "location": this.state.location
        })
      });
      // const listing = await response.json();
      // console.log(listing);
      this.props.history.push(`/listing/${listingID}`)
      console.log(this.state);
    } else {
      this.setState({ validated: false });
    }
  }
  render() {
    return (
      <div className="ListingInfo">
        <Container>
          <h1>Edit Listing</h1>
          <EditListingForm
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

export default EditListings;
