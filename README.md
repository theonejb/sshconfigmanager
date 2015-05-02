#SshConfigManager

Description
--
A utility to allow easier management of `Host` blocks in `~/.ssh/config`. I hate having to open up my editor everytime I want to add or modify options for a `Host` block in my `~/.ssh/config` file. Ubuntu has a great Gnome notification applet (or whatever it is called) to easily manage and group a list of remote hosts to SSH into. Os X didn't have anything I liked.

The plan is to create a base library that can be used to manage the config file, and then have a bunch of GUIs on top of that. The first GUI is probably going to be an `Angular` based web app, and then later maybe a Os X menu bar app.

*P.S: I'm going to create this in `Golang`, since it seems like a good practice project. I'll probably be faster in `Python`, but I **really** want to use `Golang`!*