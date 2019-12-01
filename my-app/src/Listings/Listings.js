import React, { Component } from 'react';
import {
  Card, Container, Row, Col, CardText, CardBody,
  CardTitle, CardSubtitle, Button
} from 'reactstrap';
import "./Listings.css";
import '../App.css';

class Listings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listings: []
    }
  }
  componentDidMount = async () => {
    await this.getListings();
  }
  getListings = async () => {
    // get listings here
    const listings = [];
    for (let i = 0; i < 20; i++) {
      listings.push({
        title: `Listing ${i}`,
        description: `Description ${i}`,
        condition: `Condition ${i}`,
        price: `Price: $${i}`,
        location: `Location ${i}`,
        contact: `Contact ${i}`,
        creator: `Creator ${i}`
      })
    }
    this.setState({ listings: listings });
  }
  render() {
    return (
      <div className="Listings">
        <Container>
          <h1>All Listings</h1>
          <Row>
            {this.state.listings.map(listing => {
              return (
                <Col xs="12" sm="6" md="4">
                  <Card className="listing-card">
                    {/* <CardImg top width="100%" src="/assets/318x180.svg" alt="Card image cap" /> */}
                    <CardBody>
                      <CardTitle>{listing.title}</CardTitle>
                      <CardSubtitle>{listing.price}</CardSubtitle>
                      <CardText>{listing.description}</CardText>
                      <Button>More Info</Button>
                    </CardBody>
                  </Card>
                </Col>
              );
            })}
          </Row>
        </Container>
      </div>
    );
  }
}

export default Listings;
