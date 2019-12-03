import React, { Component } from 'react';
import {
  Card, Container, Row, Col, CardText, CardBody,
  CardTitle, CardSubtitle, Button, Spinner
} from 'reactstrap';
import { Link } from 'react-router-dom';
import { array } from 'prop-types';

class MyListings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      listings: [],
      loading: false
    }
  }

  componentDidMount = async () => {
    this.setState({ loading: true })
    const userResponse = await fetch("https://api.briando.me/v1/users/me", {
      method: 'GET',
      headers: {
        'Authorization': this.props.sid
      }
    });
    const user = await userResponse.json();
    const listingsResponse = await fetch(`https://api.briando.me/v1/listings/creator/${user.id}`, {
      method: 'GET'
    });
    const listings = await listingsResponse.json();
    this.setState({ listings: listings, loading: false })
  }

  delete = async (listingID) => {
    this.setState({ loading: true })
    const response = await fetch(`https://api.briando.me/v1/listings/${listingID}`, {
      method: 'DELETE',
      headers: {
        'Authorization': this.props.sid
      }
    });
    const newListings = this.state.listings;
    for (let i = 0; i < newListings.length; i++) {
      if (newListings[i]._id === listingID) {
        newListings.splice(i, 1)
        break;
      }
    }
    this.setState({ listings: newListings })
    this.setState({ loading: false })
  }

  render() {
    if (!this.props.sid) {
      this.props.history.push("/sign-in");
    }
    return (
      <div>
        <Container>
          <h1>My Listings</h1>
          {this.state.listings.length ?
            <Row>
              {this.state.listings.map(listing => {
                const listingRoute = `/edit/${listing._id}`;
                return (
                  <Col xs="12" sm="6" md="4" key={listing._id}>
                    <Card className="listing-card">
                      <CardBody>
                        <CardTitle>{listing.title}</CardTitle>
                        <CardSubtitle>${listing.price}</CardSubtitle>
                        <CardText>{listing.description}</CardText>
                        <Link to={listingRoute}>
                          <Button>Edit</Button>
                        </Link>
                        <Button onClick={() => this.delete(listing._id)}>Delete</Button>
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
    )
  }
}

export default MyListings;