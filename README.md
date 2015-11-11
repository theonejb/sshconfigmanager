#SshConfigManager

##Update
I'm no longer working on this. I've shifted my focus to (https://github.com/theonejb/dropletconn), which I think is a better tool than what the `sshconfigmanager` could have been.

The original reason for starting this project was to make managing my SSH config easier, because I manage the cloud servers at work and I had to update my config file everytime we created/destroyed a droplet on DigitalOceans. During the course of working on this, I found a couple of other projects that did what I wanted; allow easier connectivity to my droplets. But they worked with the APIv1, which I couldn't use as I was not the admin user on our DigitalOcean account.

So, I created the `dropletconn` manager project. It allows you to list your droplets, and connect to them using just their name. No need to configure the SSH config file per server. It's been working pretty well for me, so I no longer need a SSH config manager.

Old Description ~~Description~~
--
A utility to allow easier management of `Host` blocks in `~/.ssh/config`. I hate having to open up my editor everytime I want to add or modify options for a `Host` block in my `~/.ssh/config` file. Ubuntu has a great Gnome notification applet (or whatever it is called) to easily manage and group a list of remote hosts to SSH into. Os X didn't have anything I liked.

The plan is to create a base library that can be used to manage the config file, and then have a bunch of GUIs on top of that. The first GUI is probably going to be an `Angular` based web app, and then later maybe a Os X menu bar app.

*P.S: I'm going to create this in `Golang`, since it seems like a good practice project. I'll probably be faster in `Python`, but I **really** want to use `Golang`!*
