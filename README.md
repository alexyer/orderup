# orderup
A slack bot to serve orders like a restaurant, give people an order number, cook things in serial, and mark orders complete after you serve them.

## Why?

We are busy. It makes us grumpy sometimes. And when people interrupt us on slack, we get even grumpier. We lose context, and we start to forget what we were doing last. Things fall on the floor.

We could just use a ticketing system, and tell people to file a ticket, but the culture of slack isn't like that. We like doing things in real time, and we encourage some degree of interruptions and impromptu conversations. But we do need some order to the madness, and a queue makes sense. It's better than "I'm busy, get lost." It's more like "here's what I'm up to, you're next."

So we built orderup. It's just like a restaurant. When someone asks you for something, you give them an order number, cook something up and serve it to them. Your burger will be ready after I make fries for that guy over there. Now replace burger and fries with insurance quote, git commit, phone call, whatever.

## Setup

1. Add [Slash commands](https://slack.com/apps/A0F82E8CA-slash-commands) to your slack channel(s).
2. Add a slash command as follows:
3. Command: /orderup
    URL: yourhost.com:5000/orderup
    Method: POST
    Other fields are optional.
1. Run `make build && make install`
2. `nohup orderup -host 162.243.114.162 -port 5000 -db database.db -passcode secret11722 &`

## Commands

### `/orderup help`

Shows help on all commands

### `/orderup create-q mynoodles`

This will create orders for my queue named mynoodles.

mynoodles queue created.

### `/orderup create-order mynoodles @jimuser pork sandwich`

mynoodles order 3 for @jimuser pork sandwich - order 3. There are 4 orders ahead of you.

### `/orderup finish-order mynoodles 3`

@jimuser your order is finished. Mynoodles: Order 3. Pork sandwich.

### `/orderup list mynoodles `

Mynoodles: what's cooking:

1 @jimuser - soup with chicken 

3 @jimuser - pork sandwich

24 @bethjkl - pizza burgers

### `/orderup history mynoodles`

Mynoodles history:

2 @jimuser - soup

4 @jimuser - turkey

5 @bethjkl - french fries

Etc....

### Using CURL

You can use CURL instead of Slack if you want.

`curl -H "Content-Type: application/json" -X POST -d '{"name":"mynoodles"}' http://162.243.114.162:5000/api/v1/queues`

`curl -H "Content-Type: application/json" -X DELETE -d '{"name":"mynoodles"}' http://162.243.114.162:5000/api/v1/queues`

`curl -H "Content-Type: application/json" -X POST -d '{"name":"mynoodles","user":"jimmy","description":"a hamburger"}' http://162.243.114.162:5000/api/v1/queues/order`

`curl -H "Content-Type: application/json" -X PUT -d '{"name":"mynoodles","id":"1"}' http://162.243.114.162:5000/api/v1/queues/orders/finish`

`curl -H "Content-Type: application/json" -X GET -d '{"name":"mynoodles"}' http://162.243.114.162:5000/api/v1/queues/orders/list`

`curl -H "Content-Type: application/json" -X GET -d '{"name":"mynoodles"}' http://162.243.114.162:5000/api/v1/queues/orders/history`

`curl -H "Content-Type: application/json" -X POST -d '{"name":"mynoodles","user":"jimmy","description":"a hamburger"}' http://162.243.114.162:5000/api/v1/queues/order`

If -passcode flag is specified when the service started, BasicAuth headers should be added

`curl -u :passphrase -H "Content-Type: application/json" -X POST -d '{"name":"mynoodles","user":"jimmy","description":"a hamburger"}' http://162.243.114.162:5000/api/v1/queues/order`
