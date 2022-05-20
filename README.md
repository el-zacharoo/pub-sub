# publisher

This service is apart of my `pub-sub` services, which connects to my other app, the `subscriber` service.
Pub/Sub allows services to communicate asynchronously, with latencies on the order of 100 milliseconds.
When I push data from this service, my subscriber recieves the data, puts in in mongoDB and returns a message. 
