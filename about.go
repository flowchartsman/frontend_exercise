package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const introText = `<html>
<head>
<title>Frontend Coding Exercise</title>
<link href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
<style>
.badge-get {
	background-color: #61affe;
}

.badge-post {
	background-color: #47cc90;
}

.badge-put {
	background-color: #fca131;
}

.badge-delete {
	background-color: #fa3e3f;
}

</style>
</head>
<body>
<div class="container-fluid">
<h1>Description</h1>
<p>Hello, and welcome to the frontend design exercise. It simulates a site for booking parties of different types, and is designed to provide a framework to design a frontend around and discuss problems.
<p>The exercise takes the form of a series of routes, describing a fictitious party planner app, which describes various different types of parties the user can plan, as well as the required attributes and validation for each one.
<h2>Requirements</h2>
<p>Hopefully this exercise is an intersting problem that can be solved multiple ways. Ideally we can discuss a solution that provides the following:
<ul>
<li>Some kind of rendering of each of the different types of party</li>
<li>Submission to the POST route of a valid party</li>
<li>Error Reporting</li>
</ul>

<h1>API</h1>
The application takes the form of a dummy API which is not backed by any form of datastore whatsoever. It has <code>GET</code> routes to fetch data and schema, and a <code>POST</code> route to validate a request. In order to simulate error conditions, the <code>POST</code> route will occasionally throw a <code>500</code> error in addition to normal validation.

<h2>Routes</h2>

<h3><span class="badge badge-success">GET</span><span class="text-monospace">/partytypes<span></h3>
<p>This route will return a list of valid party types you can book via the <span class="text-monospace">/bookparty</span> post routes. Note that these are case-sensitive.</p>
<div class="card text-white bg-dark mb-3">
	<div class="card-header">Example</div>
	<div class="card-body">
<pre><code class="text-white">$ curl API_URL/partytypes
["ExampleParty","AnotherParty"]</code></pre>
	</div>
</div>

<hr />
<h3><span class="badge badge-success">GET</span><span class="text-monospace">/partytype/{type name}</span></h3>
<p>This route will return the specification for the specified party type, for posting to the <span class="text-monospace">/bookparty</span> routes. The spec represents the fields of the expected datatype and some metadata about them. Let's look at an example:</p>
<div class="card text-white bg-dark mb-3">
	<div class="card-header">Example</div>
	<div class="card-body">
<pre><code class="text-white">$ curl -s API_URL/partytype/ExampleParty|python -mjson.tool
{
	"attendees": {
		"checks": [
			"gt=0"
		],
		"list": true,
		"required": true,
		"type": "string"
	},
	"end_time": {
		"checks": [
			"gtfield=start_time"
		],
		"list": false,
		"required": true,
		"type": "RFC 3339"
	},
	"start_time": {
		"checks": null,
		"list": false,
		"required": true,
		"type": "RFC 3339"
	}
}</code></pre>
	</div>
</div>
<p>Each key represents one of the JSON fields the API expects for this particular party type. Each of these has a set of keys that describe the field:</p>
<h4>required</h4>
<p>This is a boolean, and it is set to true if the field is required.</p>
<h4>list</h4>
<p>This is a boolean, and is set to true if the route expects a list/array instead of a single value.</p>
<h4>type</h4>
<p>This is a string, representing the type of the field. It can be one of the following:</p>
<table class="table">
	<thead>
		<tr>
			<th scope="col">type value</th>
			<th scope="col">description</th>
			<th scope="col">example</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td class="table-info">string</td>
			<td>a string value</td>
			<td>"hello!"</td>
		</tr>
		<tr>
			<td class="table-info">int</td>
			<td>an integer value</td>
			<td>1</td>
		</tr>
		<tr>
			<td class="table-info">RFC 3999</td>
			<td>a text value representing a timestamp conforming to the <a href="https://www.ietf.org/rfc/rfc3339.txt">RFC 3999 spec</a></td>
			<td>"2020-03-19T06:48:34+00:00"</td>
		</tr>
	</tbody>
</table>
<h4>checks</h4>
<p>This is a list of validation checks the server will enforce on the input you provide. An error will be returned by the server if any of the fields do not properly conform to these checks. Each field can have multiple checks which can be one of the following:</p>
<table class="table">
	<thead>
		<tr>
			<th scope="col">check type</th>
			<th scope="col">description</th>
			<th scope="col">example</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td class="table-info">gt</td>
			<td>for integers, the value must be greater than the supplied value. for lists, the length of the list must be at least this</td>
			<td>gt=1</td>
		</tr>
		<tr>
			<td class="table-info">gtfield</td>
			<td>compares the given field against another given field, which must be less than it.</td>
			<td>gtfield=start_time</td>
		</tr>
		<tr>
			<td class="table-info">oneof</td>
			<td>a space-separated enum list of acceptable values for a string field (case-sensitive)</td>
			<td>oneof=day night</td>
		</tr>
	</tbody>
</table>

<hr />
<h3><span class="badge badge-info">POST</span><span class="text-monospace">/bookparty</span></h3>
<p>This route will accept posts conforming to valid parties as per the spec. If it succeeds, it will return a UUID for the booking (this isn't actually stored, but might be useful to display).
<p>The route expects a json object with the following properties:</p>
<h4>party_type</h4>
This represents the type of the party you're attempting to submit. It should be one of the types returned from the <code>/partytypes</code> route.
<h4>data</h4>
This should be a JSON object conforming to the spec of the party type you're attempting to submit.
<p>Unlike <code>/bookpartyprod</code> this route will not randomly fail, so you will probably want to develop against this route, at least initially :D
<div class="card text-white bg-dark mb-3">
	<div class="card-header">Example</div>
	<div class="card-body">
<pre><code class="text-white">$ curl -XPOST API_URL/bookparty -d'{
	"party_type":"ExampleParty",
	"data": {
		"attendees":["bob smith"],
		"start_time":"2020-03-19T06:48:34+00:00",
		"end_time":"2020-03-20T06:48:34+00:00"
	}
}'
mKNkssVzbvNWRavLkCBKon</code></pre>
	</div>
</div>

<hr />
<h3><span class="badge badge-info">POST</span><span class="text-monospace">/bookpartyprod</span></h3>
<p>This route functionally identical to <code>/bookparty</code>, only it will fail 20 percent of the time, just like a real production environment! Target this route for the full "production experience" &#x1F644;.
</div>
</body>
<script type="text/javascript">
	var host = window.location.hostname;
	var port = window.location.port;
	if (port != "") {
		host = host+":"+port;
	}
	var codes = document.getElementsByTagName('code');
	[].forEach.call(codes, c => {
		c.innerText=c.innerText.replace(/API_URL/, host);
	});
</script>
</html>`

func aboutRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(introText))
}
