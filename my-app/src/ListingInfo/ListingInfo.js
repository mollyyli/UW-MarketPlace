import React, { Component } from 'react';
import {
  Card, Container, Row, Button, Spinner
} from 'reactstrap';
import { Link } from 'react-router-dom';
import "./ListingInfo.css";
import '../App.css';
import { CardBody } from 'react-bootstrap/Card';

class ListingInfo extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listing: {},
      loading: false
    }
  }
  goBack = () => {
    this.props.history.goBack();
  }
  componentDidMount = async () => {
    await this.getListingInfo();
  }
  getListingInfo = async () => {
    this.setState({ loading: true });
    const response = await fetch(`https://api.briando.me/v1/listings/${this.props.match.params.listingID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    const listing = await (response.json());
    this.setState({ listing: listing, loading: false });
  }
  render() {
    return (
      <div className="ListingInfo">
        <Container>
          {this.state.loading ?
            <div className="loading-spinner">
              <Spinner style={{ width: '8rem', height: '8rem' }} color="primary" />
            </div> :
            <div>
              <div className="listing-body">
                <h1>{this.state.listing.title}</h1>
                <div>{this.state.listing.description}</div>
                <div>{this.state.listing.contact}</div>
                <div>{this.state.listing.location}</div>
                <div>${this.state.listing.price}</div>
              </div>
              <Button onClick={this.goBack}>
                Back
              </Button>
            </div>
          }
        </Container>
      </div>
    );
  }
}

export default ListingInfo;
