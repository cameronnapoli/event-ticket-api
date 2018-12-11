# High Volume Ticket Purchasing

Application to process a high volume of requests to purchase tickets to events.

Application locks a ticket for purchase, and then allows a user to complete the ticket purchase in a time frame. If the user fails to complete the ticket purchase, then the ticket is released back into the pool.

## Dependencies

    github.com/gorilla/mux
    github.com/go-redis/redis
