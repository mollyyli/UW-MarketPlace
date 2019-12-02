import React, { Component } from 'react';
import {
  Card, Container, Row, Col, CardText, CardBody,
  CardTitle, CardSubtitle, Button, Spinner
} from 'reactstrap';
import { Link } from 'react-router-dom';
import "./ListingInfo.css";
import '../App.css';

class ListingInfo extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listing: {},
      loading: false
    }
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
    const listing = await response.json();
    console.log(listing)
    this.setState({ listings: listing, loading: false });
  }
  render() {
    return (
      <div className="ListingInfo">
        <Container>
        
        </Container>
      </div>
    );
  }
}

export default ListingInfo;
