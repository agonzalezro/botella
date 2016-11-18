Ava
===

_Ava is a temporal name since it already clashes with an important JS test runner. New names are welcomed!_

**Please, be aware that this is a beta project, don't try to run it for important stuff unless you know what you are doing!**

Slack bot with plugins. What's the selling point then? The plugins are Docker containers and the bot interacts with them via stdin/stdout. In plain English that means that you can develop your plugins in whatever language you want and you will just need to pack them as Docker container to be able to use them with Ava.

Usage
-----

You just need to be sure that you have a `ava.yaml` file in the path where you run it. This file could look like:

```yaml
adapters:
  - name: slack
    environment:
      key: xoxb-xxx
  - name: http
    environment:
      port: 8080

plugins:
  - image: agonzalezro/ava-test
    environment:
      KEY: this-is-a-secret
    only_mentions: true
```

You can easily see that we are defining a list of adapter (how to connect with the bot) and plugins that are going to be run when the bot receives a message.

### Adapters

At the moment of writing we support two types of adapters:

#### Slack

To use the slack adapter you just need to define a key in this form:

```yaml
adapters:
  - name: slack
    environment:
      key: xoxb-xxx
  ...
```

To get that key you could just go to https://your-org-here.slack.com/services/new/bot and create a new bot, after configuring it you will see an API Token, copy/paste it in your `ava.yaml`.

Probably, while you were creating the bot in Slack you saw that you could define its profile pic and name, be original!

**Note:** if for security reasons you prefer to set that API key as an environment variable you can use the environment variable `SLACK_KEY`. 

#### HTTP

The HTTP adapter is more easy to setup, you just need to define a port where you want it to be listening:

```yaml
adapters:
  - name: http
    environment:
      port: 8080 # You can also set HTTP_PORT instead setting it here
  ...
```

Now you can start POSTing, GETting or whatever to your bot in the `/` path of wherever you are running it.

For example, if we are running it in localhost in the port 8080 a good way to test it would be using [httpie](https://httpie.org/):

```bash
$ echo "hi Ava!"|http -f GET localhost:8080 Content-Type:text/plain
```

### Plugins

The plugins is just a list of docker images. Check the previous example:

```yaml
...
plugins:
  - image: agonzalezro/ava-test
    environment:
      KEY: this-is-a-secret
    only_mentions: true
```

It has 3 basic sections:

- **`image`**: it's just the name of the Docker image to be run.
- **`environment`**: the environment variables that you want to set to the container when you run it. Be careful, they are caseSensitive.
- **`only_mentions`** et al. They are basically three permissions that you can set to the plugin and that will work in the adapters that have the concept of channels and direct messages:
    - `only_mentions`: will make the bot just react when its mentioned.
    - `only_channels`: will make the bot just react when the message was send in a channel.
    - `only_direct_messages`: you know how this goes...

**Note:** If you want you can read the values of the environment keys from the host/system environment keys. Let's explain with an example:

```
...
plugins:
  - image: agonzalezro/ava-test
    environment:
      KEY:
```

In the previous configuration we say that we want our container to receive an environment variable called `KEY` but we don't set any value to it. If we want to give it one value we will just need to run ava like this:

```
AGONZALEZRO_AVA_TEST_KEY=xxx ./ava
```

Of course, how you set those variables is up to you, you don't need to do it inline as explained in the example.

Also note that all the special characters in the image name are being replaced by `_`s for compatibility reasons.

Available plugins
-----------------

| Image | Description |
| ----- | ----------- |
| [agonzalezro/ava-test](https://hub.docker.com/r/agonzalezro/ava-test/) | It's a test plugin that echoes whatever you write to it and shows the value of an environment variable called `ENV`. |


Developing plugins
------------------

Developing a plugin for Ava is extremely simple, you will need to write a program and pack it with Docker.

What are the characteristics that my program needs to follow? It will need to read a line (or several) from stdin and write a response to stdout.

In the [examples/](examples/) folder you will find an example plugin called `ava-test`. This plugin is uploaded to my Docker Hub so you can use it without building it, but it's a good place to learn.

### Examples

In the [examples/](examples/) folder you can find a simple plugin that does two important things:

- it will echo whatever message is send to it, thinking out of the box you could do whatever you want with that input (from the container stdin) and output whatever you want instead just an echo.
- it shows the value of an environment variable set on the `ava.yaml` file. What we are trying to show with this? That you can pass secrets around (from the `ava.yaml` to the container) without compromise them. Imagine that the container needs an AWS key, just add it to your `ava.yaml` and read the value from inside the container.


Developing Ava
--------------

### Compiling

You will need [Go](https://golang.org/) and [glide](https://github.com/Masterminds/glide):

```bash
$ glide install
$ go build
```

### Testing

```bash
$ go test $(glide novendor)
```
