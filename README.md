ava (temporal name :)
===

Slack bot with plugins. What's the difference? The plugins are Docker containers and the bot interacts with them via stdin/stdout.

Check [examples/](examples/) for more info.

Usage
-----

Check [examples/ava.yaml](examples/ava.yaml) for an example configuration file.

In the [examples/](examples/) folder you will find as well an example plugin called `ava-test`. This plugin is uploaded to my Docker Hub so you can use it without building it, but it's a good place to learn.

Tips
----

You will need an Slack Key to use the Slack adapter (the only one available for now).

Go to https://your-org-here.slack.com/services/new/bot and create a new bot, after configuring it you will see an API Token, you will need that one.
