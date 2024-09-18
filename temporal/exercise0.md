# Exercise 0

Time to setup Temporal on your machine!

Temporal provides a simple way to start a local server for development.

First of all, install the temporal tool:

```bash
$ go install github.com/temporalio/cli/cmd/temporal@latest
```

Test your installation:

```bash
$ temporal -h
```

If it's not working, you can download the latest release: [https://github.com/temporalio/cli/releases/tag/v1.1.0](https://github.com/temporalio/cli/releases/tag/v1.1.0).

Then try to start the development server:

```bash
$ temporal server start-dev
CLI 0.0.0-DEV (Server 1.25.0, UI 2.30.3)

Server:  localhost:7233
UI:      http://localhost:8233
Metrics: http://localhost:60293/metrics
```

You can access to the Web UI at this address: [http://localhost:8233](http://localhost:8233)

Your code will be able to contact the Temporal server with this address: `localhost:7233`

If you need to change the UI port, you can use `--ui-port 8233`.
Same for the Temporal server with `--port 7233`.

You're now ready to work on your first workflows!