# Requests

[![Join the chat at https://gitter.im/leonmaia/requests](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/leonmaia/requests?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/pkg/github.com/leonmaia/requests)
[![wercker status](https://app.wercker.com/status/93d520ff365ec9ed21189add12450999/s/master "wercker status")](https://app.wercker.com/project/bykey/93d520ff365ec9ed21189add12450999)

Most existing requests packages I've seen reimplement all features requests offers. This Request inherits all the behavior and functions of [http.Requests](https://godoc.org/net/http#Request) package and adds others functions and behaviors.

![amazing](https://raw.github.com/leonmaia/requests/master/readme_assets/jake_amazed.gif)


Features
--------

- Retries
- Connection Timeouts

Installation
------------

To install Requests, simply:

    $ go get github.com/leonmaia/requests


Usage
------------
```go
package whatever

import (
	"github.com/leonmaia/requests"
)

func GetExampleWithDefaultTimeoutAndRetries() error {
	r, err := requests.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		return err
	}

	response, err := r.Do()
	if err != nil {
		return err
	}
	// Do whatever you want with the response
	return nil
}
```
How to Contribute
------

I strongly encourage everyone whom creates a usefull custom assertion function
to contribute them and help make this package even better.

Make sure all the tests are passing, and that you have covered all the changes
you made.

