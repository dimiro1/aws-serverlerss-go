# How to use Standard HTTP handler with golang + AWS Lambda

This is very straightforward, I just need to create an http request and one http response based on the API Gateway specification.

# This is not ready production usage

It is just a small experiment to check if could be possible to execute net/http handlers inside lambda. I have plans to create a full library from this experiment, but, I do not know when it will be published :(

# Running
Install serverless framework (For easier deployment) https://serverless.com/ and then execute:

```sh
$ make deploy
```

# TODO

* Handle query strings
* Corner case for set cookie (See aws-serverless-java-core)
* Unit tests
* Documentation
* Corner cases?