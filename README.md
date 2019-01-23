# High Volume Ticket Purchasing

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Application to process a high volume of requests to purchase tickets to events.

Locks a ticket for purchase, and then allows a user to complete the ticket purchase in a time frame. If the user fails to complete the ticket purchase, then the ticket is released back into the pool.

## Dependency Commands

Run these commands to install project dependencies.

    go get github.com/gorilla/mux
    go get github.com/go-redis/redis
    go get github.com/micro/go-config


## Redis

Any Redis server can be used, but for testing, it's easiest to set up a local version. Take a look at [this link](https://redis.io/topics/quickstart) for a quick local setup on Linux. Redis configuration can be setup in `api/config.json`.

## API

Get the remaining ticket count.

    /remaining_tickets

Begin purchase of a ticket. Type can be on the strings: `GA`, `VIP`, or `ONE_DAY`.

    /buy_ticket/{type}

Complete purchase of a ticket. Token is the token received from the /buy_ticket endpoint.

    /complete_purchase/{token}


## Benchmarks

Tested with [go-wrk](https://github.com/tsliwowicz/go-wrk)

### /remaining_tickets 30 second test

Command to test
> go-wrk -c 10 -d 30 http://localhost:8000/remaining_tickets > remtick_30_1.txt


    Running 30s test @ http://localhost:8000/remaining_tickets
      10 goroutine(s) running concurrently
    60750 requests in 29.684248229s, 7.13MB read
    Requests/sec:		2046.54
    Transfer/sec:		245.82KB
    Avg Req Time:		4.886296ms
    Fastest Request:	244.287µs
    Slowest Request:	148.034965ms
    Number of Errors:	0

### /buy_ticket/{type} 30 second test

Command to test
> go-wrk -c 10 -d 30 -M POST http://localhost:8000/buy_ticket/GA > buyticket_30_1.txt

    Running 30s test @ http://localhost:8000/buy_ticket/GA
      10 goroutine(s) running concurrently
    21125 requests in 29.882631624s, 3.49MB read
    Requests/sec:		706.93
    Transfer/sec:		119.43KB
    Avg Req Time:		14.145624ms
    Fastest Request:	498.636µs
    Slowest Request:	226.50454ms
    Number of Errors:	0
