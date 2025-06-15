# **Relay Go Package**

This readme serves as a comprehensive guide to the relay Go package, explaining its interface, intended use cases, and how it enhances developer experience through flexibility.

## **Introduction**

The relay package provides a flexible and type-safe mechanism for dispatching and handling requests within your Go applications. It aims to simplify the process of defining request-response patterns and routing them to appropriate handlers, ultimately improving developer experience and code maintainability.

## **Why Use relay?**

In many applications, especially those dealing with various types of messages, commands, or events, managing the dispatch of these `requests` to the correct `handlers` can become cumbersome. relay addresses this by offering:

* **Type Safety:** By leveraging Go's generics, relay ensures that your requests and responses are strongly typed, reducing the likelihood of runtime errors and improving code readability.
* **Clear Separation of Concerns:** It promotes a clear separation between the request definition, the handler logic, and the dispatch mechanism, leading to more modular and testable code.
* **Enhanced Developer Experience:** The intuitive API allows developers to quickly define new request-response flows without boilerplate, focusing on the business logic rather than the plumbing.
* **Flexibility:** With features like custom default handlers and explicit error handling, relay adapts to various application requirements and error management strategies.
* **Centralized Request Handling:** It provides a central point for managing how different types of requests are processed, making it easier to reason about and debug your application's flow.

## **Core Concepts**

The relay package is built around a few fundamental concepts:

### **Requests (Req)**

A Request in relay is an interface that acts as a marker for your custom request types. It's designed to associate a specific request with its expected response type using Go generics.

```go
type Req[Response Res] interface {
    // This method is never called; it's purely a marker to link
    // a request type to its corresponding response type.
    Response() Response
}
```

**Example:**
```go
type EgRes struct {
    Incremented int
}

type EgReq struct {
    relay.Req[EgRes] // Links EgReq to EgRes
    Value int
}
```

By embedding `relay.Req[Response]`, you declare that EgReq is a request type and that its expected response is EgRes.

### **Responses (Res)**

Res is a simple interface representing any response type. Your custom response structs should implicitly satisfy this interface.

```go
type Res interface{}
```

**Example:**
```go
type EgRes struct {
    ExampleField int
}
```

### **Handlers (Handler)**

A Handler is a function that takes a specific Request type and returns its corresponding Response type along with an error.
```go
type Handler[Request Req[Response], Response any] func(Request) (Response, error)
```
**Example:**

```go
func EgHandler(req EgReq) (EgRes, error) {
    return EgRes{Incremented: req.Value + 1}, nil
}
```

Handlers encapsulate the business logic for processing a particular type of request.

### **The Relay Itself**

The Relay (represented by \*relay) is the core component responsible for managing and dispatching requests to their registered handlers.

* `TryNewRelay(config Config) (Relay, error)`: This constructor creates a new Relay instance. It can return `ErrInvalidConfig` if the provided configuration is not valid.

* `NewRelay(config Config) Relay`: This constructor creates a new Relay instance. It wraps `TryNewRelay` and panics when it returns error.

### **Configuration (Config)**

The Config struct is used to configure a Relay instance. It allows for specifying a DefaultHandler, which is invoked if no specific handler is found for a given request type.

* `NewConfig(defaultHandler DefaultHandler) Config`: Creates a new Config with a specified default handler.
* `NewConfigBuilder() ConfigBuilder`: Provides a fluent API for constructing Config objects.

```go
type Config struct {
    valid          bool
    defaultHandler DefaultHandler
}

type DefaultHandler func(req any) (Res, error)
```

## **Getting Started**

### **Creating a Relay Instance**

You initialize a Relay using `TryNewRelay`, typically with a Config built via NewConfigBuilder.

```go
import (
	"errors"
	"fmt"
	"github.com/ogiusek/relay"
)

func main() {
    r, err := relay.TryNewRelay(relay.NewConfigBuilder().
        DefaultHandler(func(req any) (relay.Res, error) {
            fmt.Printf("No specific handler found for request: %v\n", req)
            return nil, errors.New("unhandled request")
        }).
        Build(),
    )
    if err != nil {
        panic(fmt.Sprintf("unexpected error creating relay: %s\n", err.Error()))
    }
    // ...
}
```

### **Registering Handlers**

Handlers are registered with the Relay to associate a specific request type with its processing logic.

* `TryRegister[Request Req[Response], Response any](r *relay, handler Handler[Request, Response]) error`: Attempts to register a handler. It returns `ErrHandlerAlreadyExists` if a handler for the given request type is already registered. This is the recommended way to register handlers when you want to handle potential registration conflicts.
* `Register[Request Req[Response], Response any](r *relay, handler Handler[Request, Response])`: Registers a handler. **This function can panic** if a handler for the given request type is already registered. Use TryRegister if you want to avoid panics.


```go
type EgRes struct {
	Incremented int
}

type EgReq struct {
	relay.Req[EgRes]
	Value int
}

func EgHandler(req EgReq) (EgRes, error) {
	return EgRes{Incremented: req.Value + 1}, nil
}

// ... inside main or a similar function
relay.Register(r, EgHandler) // Panics if EgHandler is already registered for EgReq
// OR
// err := relay.TryRegister(r, EgHandler)
// if err != nil {
//     fmt.Printf("Failed to register handler: %s\n", err)
// }
```

### **Handling Requests**

Once handlers are registered, you can dispatch requests using the Handle function.

* `Handle[Request Req[Response], Response any](r *relay, request Request) (Response, error)`: Dispatches the request to the appropriate registered handler and returns the Response or an error. If no specific handler is found, the DefaultHandler (if configured) will be invoked.

```go
req := EgReq{Value: 7}
res, err := relay.Handle(r, req)
if err != nil {
    fmt.Printf("Error handling request: %s\n", err)
} else {
    fmt.Printf("Got incremented value: %d\n", res.Incremented) // Output: got incremented value: 8
}
```

## **Error Handling**

The relay package defines a few specific errors:

* `ErrInvalidConfig`: Returned by `TryNewRelay` if the provided Config is not valid (e.g., if it wasn't created using NewConfig or NewConfigBuilder).
* `ErrHandlerAlreadyExists`: Returned by TryRegister if you attempt to register a handler for a request type that already has a handler.
* `ErrHandlerNotFound`: This is the default error returned by the default DefaultHandler provided by NewConfigBuilder if no specific handler is found for a request. You can customize the default handler to return a different error or behavior.

Your custom handlers can return any error, which will be propagated back through the Handle function.

## **Advanced Usage**

### **Custom Default Handler**

The ConfigBuilder allows you to define a custom DefaultHandler. This is particularly useful for implementing fallback logic, logging unhandled requests, or returning a generic error for unknown request types.

```go
r, err := relay.TryNewRelay(relay.NewConfigBuilder().
    DefaultHandler(func(req any) (relay.Res, error) {
        fmt.Printf("Received an unhandled request of type %T: %v\n", req, req)
        // You could log this, send to a dead-letter queue, or return a specific error
        return nil, errors.New("this type of request is not supported")
    }).
    Build(),
)
if err != nil {
    panic(fmt.Sprintf("unexpected error %s\n", err.Error()))
}
```

## **Example Usage**

Here's a complete example demonstrating the usage of the relay package:

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/your-repo/relay" // Replace with your actual package path
)

// Define a custom response type
type EgRes struct {
	Incremented int
}

// Define a custom request type that links to EgRes
type EgReq struct {
	relay.Req[EgRes] // Embedding relay.Req[EgRes] establishes the link
	Value int
}

// Define a handler for EgReq
func EgHandler(req EgReq) (EgRes, error) {
	fmt.Printf("EgHandler received value: %d\n", req.Value)
	return EgRes{Incremented: req.Value + 1}, nil
}

// Another example: a request that doesn't expect a specific response (nil Res)
type PingReq struct {
	relay.Req[relay.Res] // Use relay.Res directly if no specific response struct
	Message string
}

func PingHandler(req PingReq) (relay.Res, error) {
	fmt.Printf("PingHandler received message: %s\n", req.Message)
	return nil, nil // No specific response struct, so return nil
}

func main() {
	// 1. Create a Relay instance with a custom default handler
	r, err := relay.TryNewRelay(relay.NewConfigBuilder().
		DefaultHandler(func(req any) (relay.Res, error) {
			fmt.Printf("Default handler: Unrecognized request received: %v\n", req)
			return nil, errors.New("unsupported request type")
		}).
		Build(),
	)
	if err != nil {
		panic(fmt.Sprintf("unexpected error during relay creation: %sn\", err.Error()))
	}

	// 2. Register handlers for specific request types
	relay.Register(r, EgHandler)
	relay.Register(r, PingHandler)

	// 3. Create and handle a request for which a handler is registered (EgReq)
	egRequest := EgReq{Value: 7}
	// Demonstrate JSON encoding/decoding for a request (optional, but good for demonstrating data portability)
	if bytes, err := json.Marshal(egRequest); err != nil {
		panic(fmt.Sprintf("json marshal error: %s", err.Error()))
	} else if err := json.Unmarshal(bytes, &egRequest); err != nil {
		panic(fmt.Sprintf("json unmarshal error: %s", err.Error()))
	}

	egResponse, err := relay.Handle(r, egRequest)
	if err != nil {
		panic(fmt.Sprintf("error handling EgRequest: %s\n", err))
	}
	fmt.Printf("Handled EgRequest: got incremented value %d\n", egResponse.Incremented)

	// 4. Create and handle another registered request (PingReq)
	pingRequest := PingReq{Message: "Hello, Relay!"}
	_, err = relay.Handle(r, pingRequest) // PingHandler returns nil for response
	if err != nil {
		panic(fmt.Sprintf("error handling PingRequest: %s\n", err))
	}
	fmt.Println("Handled PingRequest successfully.\n")

	// 5. Create and handle an unregistered request to trigger the default handler
	type UnregisteredReq struct {
		relay.Req[relay.Res]
		Data string
	}
	unregisteredRequest := UnregisteredReq{Data: "Some unknown data"}
	_, err = relay.Handle(r, unregisteredRequest)
	if err != nil {
		fmt.Printf("Handled unregistered request, as expected: %s\n", err.Error())
	}
}
```

## **Contributing**

Contact us.

## **license**

This isn't something i expect to happen therefor i do not has policy for this yet
