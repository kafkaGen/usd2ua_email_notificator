# Currency Mail Service

This project implements a service with an API for managing email subscribers, retrieving current currency rates, and sending daily email messages to subscribed users. The service includes robust error handling, validation, and scheduled tasks using a cron job.

## Features
1. Email Subscription Management:

    - Add Subscriber: API to add new email subscribers.
    - Delete Subscriber: API to remove existing email subscribers.
    - List Subscribers: API to retrieve the list of all email subscribers.

2. Currency Rate Retrieval:

    - Get Current Rate: API to get the current currency rate.

3. Email Validation and Error Handling:

    - Validate email format.
    - Prevent adding duplicate email addresses.
    - Handle attempts to remove non-existent email addresses.
    - Handle various runtime errors gracefully.

4. Daily Email Notifications:

    - Send daily email messages to all subscribed users using Google SMTP server.
    - Scheduled task with a cron job to update each subscriber with new emails daily.

5. Docker and Docker Compose:

    - Separate Docker images for the API server and the cron service (to let each service response only for one task)
    - PostgreSQL database managed by Docker Compose.
    - On PostgreSQL database startup triggered init scripts for setting up DB tabels.
    - Two-stage Docker build for optimized image size, using a lightweight Alpine image for the final stage.
    - Sensitive variables stored in a .env file (example file provided without actual keys/credentials).

## Technologies
    - Language: Go (version 1.18)
    - Database: PostgreSQL
    - Email: Google SMTP server
    - Containerization: Docker and Docker Compose

## Setup and Run
1. Create and configure .env file:

    Copy .env.example to .env
    Fill in the required environment variables (Database credentials, SMTP credentials, API key, etc.)

2. Build and Start Containers:
    ```bash
    docker compose --env-file .env -f docker/docker-compose.yaml up -d --build
    ```

3. Access the API:

    The API server will be accessible at http://localhost:8080

4. Endpoints:

 - Add Subscriber: POST /subscribe \
    Request Body: {"email": "user@example.com"}
    ```
    curl -X POST http://localhost:8080/subscribe \
     -H "Content-Type: application/json" \
     -d '{"email":"example@example.com"}'
    ```

- Delete Subscriber: POST /unsubscribe \
    Request Body: {"email": "user@example.com"}
    ```
    curl -X POST http://localhost:8080/unsubscribe \
     -H "Content-Type: application/json" \
     -d '{"email":"example@example.com"}'
    ```

- List Subscribers: GET /subscribers
    ```
    curl -X GET http://localhost:8080/subscribes
    ```

- Get Current Rate: GET /rate
    ```
    curl -X GET http://localhost:8080/rate
    ```

## Project Logic
The PostgreSQL database table Subscribers is used to store the emails of users. Three endpoints are developed to handle subscriptions: add one, delete one, and get the list of subscribers. To automate database setup, when starting the database with Docker Compose, table initialization scripts are sent to the docker-entrypoint-initdb.d directory. This ensures that all necessary tables are created at startup.

The project includes a rate endpoint to get the current USD to UAH currency rate. Under the hood, this endpoint performs an API call to the Exchange Rate API.

The cron package in Go is used to schedule the service that sends emails to users. The service is configured to run daily using the @daily cron setting. The project contains a message template in the messages folder. The cron job retrieves this template, fetches the current currency rate, and sends the message to all emails in the database using the SMTP service provided by Google. For this, app credentials are created in a Google account.

## Author Feedback
For me, it was a very exciting task and a great first experience with Go language. I find Go to be a very convenient language that deserves further investigation. Setting up the cron task was a good challenge, as logging and error handling are complicated, and it was hard to make everything work together.