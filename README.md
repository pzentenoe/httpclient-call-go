# httpclient-call-go

## Overview

The `httpclient-call-go` library simplifies making HTTP calls to various API services efficiently and straightforwardly.
It is designed to seamlessly integrate into any Go project requiring HTTP API interactions.

![CI](https://github.com/pzentenoe/httpclient-call-go/actions/workflows/actions.yml/badge.svg)
![Quality Gate](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=httpclient-call-go&metric=alert_status&token=sqb_28f943efa72bc60b8e1c5447065df406ec45ef08)
![Coverage](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=httpclient-call-go&metric=coverage&token=sqb_28f943efa72bc60b8e1c5447065df406ec45ef08)
![Bugs](https://sonarqube.vikingcode.cl/api/project_badges/measure?project=httpclient-call-go&metric=bugs&token=sqb_28f943efa72bc60b8e1c5447065df406ec45ef08)

### Buy Me a Coffee

<a href="https://www.buymeacoffee.com/pzentenoe" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

Thank you for your support! ❤️

## Features

- Easy HTTP client configuration.
- Full support for customizing HTTP requests (headers, body, timeouts).
- Convenient methods for making HTTP calls and deserializing JSON responses.

## Installation

To use `httpclient-call-go` in your project, install it using the following Go command:

```bash
go get github.com/pzentenoe/httpclient-call-go
```

## Quick Start

### Setting Up HTTP Client

First, import the library and create a new instance of HTTPClientCall specifying the base URL of the API service and an
HTTP client:

```go
import (
"net/http"
client "github.com/pzentenoe/httpclient-call-go"
)

httpClientCall := client.NewHTTPClientCall("https://dummyhost.cl", &http.Client{})
```

## Making an HTTP Call

Using **Do** Implementation

To perform a simple POST request and handle the response as []byte:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	client "github.com/pzentenoe/httpclient-call-go"
)

func main() {
	httpClientCall := client.NewHTTPClientCall("https://dummyhost.cl", &http.Client{})
	headers := http.Header{
		client.HeaderContentType: []string{client.MIMEApplicationJSON},
	}
	dummyBody := map[string]interface{}{"age": 30, "name": "test"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := httpClientCall.
		Method(http.MethodPost).
		Path("/path").
		Body(dummyBody).
		Headers(headers).
		Do(ctx)
	if err != nil {
		fmt.Println("Error calling the API:", err)
		return
	}
	defer response.Body.Close()

	dataBytes, errToRead := io.ReadAll(response.Body)
	if errToRead != nil {
		fmt.Println("Error reading data:", errToRead)
		return
	}
	fmt.Println(string(dataBytes))
}
```

Using **DoWithUnmarshal** Implementation

To perform a POST request and automatically deserialize the JSON response into a Go structure:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	client "github.com/pzentenoe/httpclient-call-go"
)

type someBodyResponse struct {
	Name string `json:"name"`
}

func main() {
	httpClientCall := client.NewHTTPClientCall("https://dummyhost.cl", &http.Client{})
	headers := http.Header{
		client.HeaderContentType: []string{client.MIMEApplicationJSON},
	}
	dummyBody := map[string]interface{}{"age": 30, "name": "test"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var responseBody someBodyResponse
	resp, err := httpClientCall.
		Method(http.MethodPost).
		Path("/path").
		Body(dummyBody).
		Headers(headers).
		DoWithUnmarshal(ctx, &responseBody)
	if err != nil {
		fmt.Println("Error calling the API:", err)
		return
	}
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Name in Response:", responseBody.Name)
}
```

## Testing

Execute the tests with:

```bash
go test ./...
```

## Contributing
We welcome contributions! Please fork the project and submit pull requests to the `main` branch. Make sure to add tests
for new functionalities and document any significant changes.

## License
This project is released under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Changelog
For a detailed changelog, refer to [CHANGELOG.md](CHANGELOG.md).

## Autor
- **Pablo Zenteno** - _Full Stack Developer_ - [pzentenoe](https://github.com/pzentenoe)
