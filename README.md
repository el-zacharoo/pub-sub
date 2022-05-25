# publisher

This service is apart of my `pub-sub` services, which connects to my other app, the `subscriber` service.
Pub/Sub allows services to communicate asynchronously, with latencies on the order of 100 milliseconds.
When I push data from this service, my subscriber recieves the data, puts in in mongoDB and returns a message. 

## Get started 
- Install/update dependencies `go get -u`
- Run `make proto`, this will also send the proto definitions to `https://buf.build` as a commit. 
- Run the application: `make run`, this will run on port 8081
- Make sure the subscriber service is also running, this should be visibil on port 8082

### Dapr ###

Install and configure the dapr CLI:

```bash
brew install dapr/tap/dapr-cli
dapr init
```

### Buf ###

Add Buf API key to `.netrc` file

```bash
machine buf.build password <your Buf API key>
machine go.buf.build login <your Buf username> password <your Buf API key>
```
## Postman Setup ##
- Go to New, select gRPC Request
- Call `localhost:8080`
- Select the proto definitions file to put into the methods input
- Make sure server reflection is enabled in methods
- You can generate example JSON Messages (this is particularly useful when using the CREATE method, make sure the ID field isn't present as this is being     filled in automatically when an entry is created) 
