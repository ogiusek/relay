# `relay`

`github.com/ogiusek/relay`

The `relay` package provides a robust and flexible mechanism for decoupling communication concerns within your application. By enforcing upfront registration of handlers for specific request types, `relay` promotes good software design principles, enabling you to easily swap communication methods (e.g., HTTP, TCP, or even in-memory calls for backend/frontend in a single executable) without altering core business logic.

## Features

  * **Decoupled Communication:** Abstract away the underlying communication mechanism.
  * **Strongly Typed Handlers:** Ensure type safety between requests and responses.
  * **Explicit Registration:** Forces clear definition of how requests are handled.
  * **Default Handler:** Provides a fallback for unhandled request types.
  * **Composable Builder:** Allows for flexible and chained configuration.
  * **Middlewares:** Allows registering middlewares (`Builder.RegisterMiddleware`)

## Installation

To use `relay`, simply run:

```bash
go get github.com/ogiusek/relay
```

## Core Concepts

At the heart of the `relay` package are two main components:

  * **`Builder`**: Used to construct and configure a `Relay` instance. It provides methods for registering handlers and setting a default handler.
  * **`Relay`**: The operational component that dispatches requests to their registered handlers.

### Requests and Responses

The `relay` package uses a generic approach for requests and responses. Your request types must embed `relay.Req[Response]`, which serves as a marker to bind a request to its expected response type.

```go
type MyRequest struct {
    relay.Req[MyResponse] // Binds MyRequest to MyResponse
    // Your request fields here
}

type MyResponse struct {
    // Your response fields here
}
```

## Usage

Building and using a `Relay` typically involves these steps:

1.  **Define your Request and Response types.**
2.  **Create Handler Functions** for each request type.
3.  **Use `NewBuilder()`** to start configuring your `Relay`.
4.  **Register Handlers** using `Register`.
5.  **Optionally set a `DefaultHandler`** for unhandled requests.
6.  **Call `Build()`** to get your `Relay` instance.
7.  **Handle requests** using the `Handle` method of the `Relay`.

### Example

Let's illustrate with a simple example:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/ogiusek/relay"
)

// 1. Define Request and Response types
type EgReq struct{ relay.Req[int] }

// Custom error for demonstration
var ErrCustomError = errors.New("a custom error occurred")

func main() {
	// 2. Create a handler function
	handler := func(req EgReq) (int, error) {
		fmt.Printf("Received request: %+v\n", req)
		return 123, nil // Return some integer response
	}

	// 3-6. Build the Relay
	r := relay.NewBuilder().
		// Using Wrap for cleaner registration, especially with generics
		Wrap(func(b relay.Builder) relay.Builder {
			return relay.Register(b, handler)
		}).
		// Set a default handler for requests without a specific handler
		DefaultHandler(func(ctx relay.AnyContext) {
			fmt.Printf("No handler found for request type: %T\n", ctx.Req())
			ctx.SetErr(ErrCustomError)
		}).
		Build()

	// 7. Handle a request
	payload := EgReq{} // Create an instance of your request
	res, err := relay.Handle[EgReq, int](r, payload)

	if err != nil {
		fmt.Printf("Error handling request: %v\n", err)
	} else {
		fmt.Printf("Response: %d\n", res) // Output: Response: 123
	}

	// Example of an unhandled request
	type UnhandledReq struct{ relay.Req[string] }
	_, err = relay.Handle[UnhandledReq, string](r, UnhandledReq{})
	if err != nil {
		fmt.Printf("Error handling unhandled request: %v\n", err) // Output: Error handling unhandled request: a custom error occurred
	}
}
```

### Explanation of Builder Methods

  * **`NewBuilder()`**: The constructor for `Builder`. It initializes the internal `relay` structure with an empty set of handlers and a default handler that returns `ErrHandlerNotFound`.

    ```go
    builder := relay.NewBuilder()
    ```

  * **`Wrap(wrapped func(Builder) Builder) Builder`**: This method allows for a more functional and composable way to extend the `Builder`'s configuration. It takes a function that accepts a `Builder` and returns a modified `Builder`. This is particularly useful for applying multiple registrations or complex configurations.

    ```go
    builder := relay.NewBuilder().
        Wrap(func(b relay.Builder) relay.Builder {
            return relay.Register(b, myHandler1)
        }).
        Wrap(func(b relay.Builder) relay.Builder {
            return relay.Register(b, myHandler2)
        })
    ```

  * **`DefaultHandler(handler DefaultHandler) Builder`**: Sets a global fallback handler for any request type that does not have a specific handler registered. If no `DefaultHandler` is set, the built `Relay` will use an internal default that returns `ErrHandlerNotFound`.

    ```go
    builder := relay.NewBuilder().
		DefaultHandler(func(req any) (relay.Res, error) {
			fmt.Printf("No handler found for request type: %T\n", req)
			return nil, ErrCustomError
		}).
    ```

  * **`Register[Request Req[Response], Response any](b Builder, handler Handler[Request, Response]) Builder`**: Registers a handler for a specific request type. This function panics if:

      * The `Builder` was not created using `NewBuilder()`.
      * A handler for the given `Request` type has already been registered.
        For non-panicking registration, consider implementing a `TryRegister` function within your application if needed.

    ```go
    type MyRequest struct{ relay.Req[string] }
    myHandler := func(req MyRequest) (string, error) {
        return "hello", nil
    }
    builder := relay.NewBuilder()
    builder = relay.Register(builder, myHandler) // Directly register
    ```

  * **`Build() Relay`**: Finalizes the `Builder` configuration and returns a `Relay` instance ready to handle requests. This method also panics if the `Builder` was not created via `NewBuilder()`.

    ```go
    r := relay.NewBuilder().Build()
    ```

## Error Handling

The `relay` package defines the following errors:

  * **`ErrDidntUseCtor`**: Indicates that a `Builder` method was called on an uninitialized `Builder` instance (i.e., `NewBuilder()` was not used).
  * **`ErrHandlerAlreadyExists`**: Signifies an attempt to register a handler for a request type that already has a handler.
  * **`ErrHandlerNotFound`**: Returned by the default handler when no specific handler is found for a given request.
  * **`ErrInvalidType`**: Returned by to middleware when wrong type is set.

It's important to note that `Register` and `Build` methods can `panic` under certain conditions. This design choice forces the client to define the `Relay` upfront and correctly, ensuring a predictable and stable state at runtime.

## Contributing

Contributions are welcome\! Please feel free to open issues or submit pull requests on the [GitHub repository](https://www.google.com/search?q=https://github.com/ogiusek/relay).
