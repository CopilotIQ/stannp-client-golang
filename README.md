# Stannp API Client

This is a Go client for the Stannp API, which allows you to send letters using the Stannp service. This README provides instructions on how to initialize a new instance of the Stannp client and how to call the `SendLetter` method to send letters.

## Prerequisites

Before using this client, make sure you have the following:

- Go programming language installed: [Go Installation Guide](https://golang.org/doc/install)

## Installation

To use the Stannp API client in your Go project, follow these steps:

1. Open a terminal or command prompt.

2. Navigate to your Go project directory.

3. Run the following command to add the client as a dependency:

   ```bash
   go get github.com/CopilotIQ/stannp-client-golang


## Import the client in your Go code
`import "github.com/CopilotIQ/stannp-client-golang/stannp"`

## Usage

Initializing a Stannp Client

To initialize a new instance of the Stannp client, you can use the New function. The New function allows you to configure the client with various options using functional options.

Here's an example of initializing a Stannp client:
```
    api := stannp.New(
    stannp.WithAPIKey("your-api-key"),
    stannp.WithTest(true),
    stannp.WithClearZone(true),
    stannp.WithDuplex(true),
    stannp.WithPostUnverified(false),
    )
```

Replace "your-api-key" with your Stannp API key. You can obtain an API key by signing up on the Stannp website and generating an API key from your account settings.

You can customize other options based on your requirements. Refer to the method signatures and documentation for more details.


## Sending a Letter

To send a letter using the Stannp client, you can call the SendLetter method. This method takes a *letter.Request object as input and returns a *letter.Response object and a *util.APIError object.

Here's an example of calling the SendLetter method:

```
recipient := letter.RecipientDetails{
    Address1:  "123 Main St",
    Country:   "United States",
    Firstname: "John",
    Lastname:  "Doe",
    State:     "California",
    Title:     "Mr",
    Town:      "Los Angeles",
    Zipcode:   "90001",
}

mergeVariables := letter.MergeVariables{
    "variable1": "value1",
    "variable2": "value2",
}

request := &letter.Request{
    MergeVariables: mergeVariables,
    Recipient:      recipient,
    Template:       "template-name",
}

response, err := api.SendLetter(request)
if err != nil {
    // Handle error
}

// Access the response data
fmt.Println("Letter ID:", response.Data.Id)
fmt.Println("Letter Status:", response.Data.Status)
```

Replace "template-name" with the name of the template you want to use for the letter. Make sure to provide appropriate recipient details and merge variables based on your letter content.

Check the letter.Request and letter.Response structures for all available fields and customize them as needed.


## Examples

For more usage examples, refer to the examples provided in the examples directory of this repository.

## Contributing

If you want to contribute to this project, feel free to submit pull requests or open issues on the GitHub repository.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
