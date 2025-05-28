## Setup Steps

### 1. Clone the Repository

First, clone the project repository to your local machine:

```bash
git clone https://github.com/nikv1811/movierental.git
cd movierental
```

### 2. Database Setup

This application uses PostgreSQL.

1.  **Create a PostgreSQL database**:
    You can create a new database using `psql` (PostgreSQL command-line client) or a GUI tool like DBeaver or pgAdmin.

    ```bash
    psql -U your_username -c "CREATE DATABASE movie-rental-db;"
    ```

    Replace `your_username` with your PostgreSQL username.

2.  **Ensure you have a user with access to this database.**

### 3. Environment Variables

Create a `config.json` file in the root directory of the project (e.g., `movie-rental-app/config.json`). This file will hold your application's configuration, including database credentials and external API keys.

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "your_db_user",
    "password": "your_db_password",
    "dbname": "movie-rental-db",
    "sslmode": "disable"
  },
  "movie_api": {
    "base_url": "https://movie-database-api1.p.rapidapi.com",
    "headers": {
      "rapid_api_host": "movie-database-api1.p.rapidapi.com",
      "rapid_api_key": "YOUR_RAPIDAPI_KEY"
    }
  }
}
```

### 4. Install Dependencies

Navigate to the project root directory in your terminal and install the Go modules:

```bash
go mod tidy
```

### 5. Run Database Migrations

The application uses GORM for ORM. The `User` and `Cart` tables will be created automatically upon the first run of the application if they do not already exist, based on the models defined in `pkg/models/user.go` and `pkg/models/requests/cart.go` (implicitly, as `requests.Cart` is used to define the cart structure).

for any further migrations required run the

```bash
go run migration/migration.go
```

### 6. Run the Application

To start the Go application:

```bash
go run main.go
```

You should see output similar to this, indicating the database connection is successful and the server is starting:

Connected to database
Server is running on port 8080
The application will be running on http://localhost:8080.

## API Documentation (Swagger)

Once the application is running, you can access the API documentation (generated using Swagger) by navigating to:

`http://localhost:8080/docs/index.html`

This page provides an interactive interface to explore and test all available API endpoints.
