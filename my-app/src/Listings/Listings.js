import React, { Component } from 'react';
import {
  Card, Container, Row, Col, CardText, CardBody,
  CardTitle, CardSubtitle, Button, Spinner
} from 'reactstrap';
import "./Listings.css";
import '../App.css';

class Listings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listings: [],
      loading: false
    }
  }
  componentDidMount = async () => {
    await this.getListings();
  }
  getListings = async () => {
    this.setState({ loading: true });
    const response = await fetch('https://api.briando.me/v1/listings', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    const listings = await response.json();
    this.setState({ listings: listings, loading: false });
  }
  render() {
    return (
      <div className="Listings">
        <Container>
          <h1>All Listings</h1>
          {this.state.listings.length ?
            <Row>
              {this.state.listings.map(listing => {
                return (
                  <Col xs="12" sm="6" md="4" key={listing._id}>
                    <Card className="listing-card">
                      <CardBody>
                        <CardTitle>{listing.title}</CardTitle>
                        <CardSubtitle>${listing.price}</CardSubtitle>
                        <CardText>{listing.description}</CardText>
                        <Button>More Info</Button>
                      </CardBody>
                    </Card>
                  </Col>
                );
              })}
            </Row> :
            this.state.loading ?
              <div className="loading-spinner">
                <Spinner style={{ width: '8rem', height: '8rem' }} color="primary" />
              </div> :
              <p>No listings found.</p>

          }
        </Container>
      </div>
    );
  }
}

export default Listings;
