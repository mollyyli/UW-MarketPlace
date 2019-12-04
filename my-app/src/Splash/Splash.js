import React, { Component } from 'react';
import { Container, Alert } from 'reactstrap';
import './Splash.css';

class Splash extends Component {
  constructor(props) {
    super(props);
    this.state = {
    }
  }
  render() {
    return (
      <div id="outsidediv">
            <div id="reviews">
                <div class="row">
                    <h2 class="col-sm">
                    Innovative.
                    </h2>
                    <h2 class="col-sm">
                    The Craigslist killer.
                    </h2>
                    <h2 class="col-sm">
                    This is peak technology.
                    </h2>
                </div>
                <div class="row">
                    <h2 class="col-sm">
                        -New York Times
                    </h2>
                    <h2 class="col-sm">
                        -The Wall Street Journal
                    </h2>
                    <h2 class="col-sm">
                        -Washington Post
                    </h2>
                </div>
            </div>
            <h1 id="headerSplash">UW MarketPlace</h1>
            <img id="img" src={require('../uwmarketplace.jpg')}></img>
        </div>
    );
  }
}

export default Splash;
