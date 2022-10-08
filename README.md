## httpclient-call-go

This library is used to make http calls to different API services

## Install Package

`go get github.com/pzentenoe/httpclient-call-go`

# Quick Start

The following is the minimum needed code to call api with httpclient-call-go

### Using Do Implementation

```go
package main

import (
	"fmt"
	"io"
	"net/http"

	client "github.com/pzentenoe/httpclient-call-go"
)

func main() {

	httpClientCall := client.NewHTTPClientCall(&http.Client{}).Host("https://dummyhost.cl")

	headers := http.Header{
		client.HeaderContentType: []string{client.MIMEApplicationJSON},
	}
	dummyBody := make(map[string]interface{})
	dummyBody["age"] = 30
	dummyBody["name"] = "test"

	response, err := httpClientCall.
		Method(http.MethodPost).
		Path("/path").
		Body(dummyBody).
		Headers(headers).
		Do()

	if err != nil {
		fmt.Println("Error to call api")
		return
	}
	//Close body response 
	defer response.Body.Close()

	datBytes, errToRead := io.ReadAll(response.Body)
	if errToRead != nil {
		fmt.Println("Error to read data")
		return
	}
	fmt.Println(string(datBytes))

}
```

### Using DoWithUnmarshal Implementation

```go
package main

import (
	"fmt"
	http "net/http"

	client "github.com/pzentenoe/httpclient-call-go"
)

type someBodyResponse struct {
	Name string `json:"name"`
}

func main() {

	httpClientCall := client.NewHTTPClientCall(&http.Client{}).Host("https://dummyhost.cl")

	headers := http.Header{
		client.HeaderContentType: []string{client.MIMEApplicationJSON},
	}
	dummyBody := make(map[string]interface{})
	dummyBody["age"] = 30
	dummyBody["name"] = "test"

	var responseBody *someBodyResponse
	resp, err := httpClientCall.
		Method(http.MethodPost).
		Path("/path").
		Body(dummyBody).
		Headers(headers).
		DoWithUnmarshal(&responseBody)

	if err != nil {
		fmt.Println("Error to call api")
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(responseBody.Name)

}
```