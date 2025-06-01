# E-Commerce Demo Boilerplate

This is a ready-to-use boilerplate project made with Go. It follows a clean and organized code structure (called clean architecture) and is good for learning or starting real applications.

This demo helps you:

- See how _products_ are grouped into _categories_ (like bread under bakery).

- Create _orders_ as a logged-in _customer_.

- Send a text message to the _customer_ to confirm the _order_.

## How It Works
### Products and Categories

- Products (like "bread") belong to categories (like "bakery").

- Categories can be nested (a category can have a sub-category).

### Customers

- A customer has a name, email, phone, and a unique ID.

### Orders

- A customer can create an order that includes one or more products.

- The order is saved in the database.

- A confirmation SMS is sent to the customer.

## Technologies Used
- **Docker** – Runs the app in containers so it works the same everywhere

- **Gqlgen** – A tool for building GraphQL APIs in Go

- **pgx** – A Go library for working with PostgreSQL databases

- **Africa's Talking** – Used to send SMS messages

- **PostgreSQL** – The database used to store everything

- **Auth0** – Handles authentication (login and user identity)

## Project Setup

#### 1. Clone the project:

`git clone git@github.com:godfreyowidi/simple-ecomm-demo.git`\
`cd simple-ecomm-demo`

#### 2. Set up environment variables:

Create a .env file or export these in your terminal:

`DATABASE_URL=postgres://postgres:password@localhost:5432/simple_ecomm?sslmode=disable`\
`AFRICASTALKING_API_KEY=your_api_key`\
`AFRICASTALKING_USERNAME=your_username`

`AUTH0_DOMAIN=your-auth0-domain`\
`AUTH0_AUDIENCE=your-auth0-api-audience`

For testing:

`TEST_DATABASE_URL=postgres://postgres:password@localhost:5432/test_ecomm?sslmode=disable`

This project uses Auth0 to protect some GraphQL operations (like placing an order). 

When a customer is authenticated, they get a JWT token.
This token proves who they are.

You must send this token in the Authorization header when making a GraphQL request like in postman or graphql playground

`Authorization: Bearer <your-auth0-token>`

#### 3. How the backend uses Auth0
- The backend checks the token with Auth0.

- If the token is valid, it pulls the user info (like `sub` which is a unique user ID).

- It uses that `sub` to link the user to a customer in the database.

- If the token is missing or invalid, protected GraphQL operations will fail.

## Running the app
Start with docker

`docker-compose up --build`

GraphQL Playground will be available at:

`http://localhost:8080`

You can test authenticated queries by clicking "HTTP HEADERS" and adding:

`
{
  "Authorization": "Bearer YOUR_JWT_TOKEN"
}
`

## Running tests

`go test -v ./tests/integration -run ^TestCreateOrderMutation$`

Make sure your __TEST_DATABASE_URL__ is set.

## Known bugs

## TODOs

## License




