# UW MarketPlace
Brian Do, Molly Li, Justin La

## Project Description

### Target Audience
UWMP is a listing application where UW students or residents within the UDistrict can participate in buy/sell exchanges. The users will be able to post listings in those categories with a specific item type.
### Users
Users may want to use this application because it is an online base where individuals within the same area can sell and buy items locally. In specific scenarios, it can be used to sell furniture when moving out, or getting rid of old clothes that are still in wearable condition. Potential users may want to sell or buy items at a discounted price, look for items that are not available in retail, or make some extra cash. The users will be able to browse selections, save as favorites, and create listings. They do this by creating creating an account with simple identifying information. Within a user account, they then have the option to create a listing. Users who do not have an account are unable to post a listing. 
### Developers
As developers we want to create this application because we think it would be a great addition to student life. There is rarely a point where people are not interested in buying/selling/trading things that they have or want, and this would be a platform to mediate those exchanges. The actual transaction is not done over the application and that is attributed to the idea of being local listing application. 

## Technical Description
### Architectural Diagram
![chart](chart.png)



### Endpoints
* GET v1/listings. Retrieves list of all current listings for sale
* GET v1/listings/{ID} Retrieves all listings of a user
* POST v1/listings/{ID}. Adds a listing to the database with given ID
* PATCH v1/listings/{ID}. Edits listing at given ID
* POST v1/users. Adds user to database
* PATCH v1/users. Edit user data, adding favorite listings
* DELETE v1/listings/{ID}. Deletes listing at given ID
### Database schemas
Listing : {
	ID: string
	Title: string
	Description: string
	Condition: string
	Price: string
	Location: string
	Contact: string
	Creator: string
}

User : {
	ID: string
	UserName: string
	Favorites: Listing[]
}


![table](table.png)




### User Stories
* P0: To get the current listings, a get request to the server would be made in order to send a response that contains the current listings

* P1: We would authenticate without external high level libraries and instead implement user authentication from scratch. To make a listing, there would be a POST request to the server that adds it to the database. 

* P2: After authenticating, to add favorites to a user, we would make a POST request on the listing

* P3: The seller can then PATCH to edit the listing or DELETE to remove the listing 
