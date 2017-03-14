[![Build Status](https://travis-ci.org/HarryEMartland/orderly-badger.svg?branch=master)](https://travis-ci.org/HarryEMartland/orderly-badger)

# Orderly Badger

Orderly Badger is a simple docker management tool designed to aid shutting down short lived testing environments. 
To stop a container after a given amount of time add an environment variable called `MAX_AGE` to the container with the desired duration to shut it down after.

Accepted durations include seconds (s), minutes (m) and hours (h) for example to stop a container after 2 hours set `MAX_AGE=2h`.
A full docker command may look like `docker run -d -e MAX_AGE=2h httpd`.

A simple web ui is included to monitor which containers will be automatically deleted. 
The container created time along with the max age and how long till the container will be stopped are show.
The ui can be accessed on port 8080.

The application listens on the docker socket to get notified when new containers are created and if they have `MAGE_AGE` set.
To start up the application and listen on the socket the socket must be mounted similar to below.

`docker run -v -d /var/run/docker.sock:/var/run/docker.sock -p8080:8080 harrymartland/orderly-badger`