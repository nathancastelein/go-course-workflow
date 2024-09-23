# Exercise 2

Let's now improve a bit our helloworld and add a step to our process.

Professor Oak is a bit sad and want to say hello to you!

## A new activity

Add a new activity, named `SayHelloToProfessorOak`.

Here is the code to say hello:

```go
    resp, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

    return result, nil
```

Add this code as the body of your activity.

Then add the call to the activity inside your Helloworld workflow.

## Run the new workflow

Before running the workflow, don't forget to register your activity in your worker!

Then run the worker, run the starter.

Looks like it's not working...

Open Temporal's UI. Open your running workflow. There's an error!

## Running the API

Run the [API](./cmd/api/main.go) to contact Professor Oak.

Observe your workflow through Temporal's UI. It should resolve by itself!

## Bonus: fix the tests

Fix the workflow test, as it's broken for now!