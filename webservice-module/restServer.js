
var http = require('http');
var fs = require('fs');
var url = require('url');
var chaincode_invoke = require('./invoke.js');//should be changed based on folder where invoke.js is kept
var chaincode_query = require('./query.js');//should be changed based on folder where query.js is kept

//Handle HTTP requests
async function handle_incoming_request (req, res) {

    console.log("INCOMING REQUEST: " + req.method + " " + req.url);
	
	 if(req.method == 'POST') {
		console.log("POST");
		switch(req.url){
			case '/sensorreading':
				var receivedReading = '';
				req.on('data', function(chunk) {
					receivedReading += chunk;
					console.log("Received data chunk " + receivedReading);
					const data = JSON.parse(receivedReading);
					chaincode_invoke.createRecord(data.temp);
				});
				req.on('end', function() {	
					console.log("Received Data: " + receivedReading);
				});
				res.writeHead(200, {'Content-Type': 'text/html'});
				res.end('post received');
				break;
			default:
				res.writeHead(404, "Not found", {'Content-Type': 'text/html'});
				res.end('Not found');
				console.log("[404] " + req.method + " to " + req.url);
				break;
		}
    }
	
	
	if(req.method == 'GET')
	{
		console.log("GET");
		switch(req.url){
			case '/sensor/all':
				let result = await chaincode_query.queryRecord();
				res.writeHead(200, {'Content-Type': 'application/json'});
				res.end(result);
				break;
			default:
				res.writeHead(404, "Not found", {'Content-Type': 'text/html'});
				res.end('Not found');
				console.log("[404] " + req.method + " to " + req.url);
				break;
		}
	
	}
}

//Create http server
var s = http.createServer(handle_incoming_request);
console.log("HTTP Server created");
s.listen(8080);
console.log("listening on port 8080");
console.log("HTTP Server created");

