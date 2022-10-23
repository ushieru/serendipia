# Serendipia
<p align="center">
  <img src="./resources/serendipia.png" alt="Serendipia" />
</p>

Serendipia is a simple RESTful (Representational State Transfer) gateway service for the purpose of discovery, load balancing and failover.

## Documentation

#### Quickstart
```bash
$ serendipia
```
## How works?
#### Register MicroService
```js
// Express example
import express from 'express'
import fetch from 'node-fetch';

// ... [code]

app.listen(0, function () {
    const port = this.address().port
    const register = () => axios.post('http://127.0.0.1:3000/services', {
            service_name: "example_service",
            service_port: port.toString()
        });

    // Register microservice 
    register()
    // Update microservice 
    setInterval(register, 4000)

    console.log(`Server run: http://localhost:${port}/`)
})
```

#### Now you only need make a request with this rules:

**[SerendipiaServer]**/**[MicroServiceName]**/**[MicroServicePath]**

## In progress
- [ ] Version support
