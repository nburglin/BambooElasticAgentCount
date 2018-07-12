## Description

This is a simple Go program designed to leverage PRTG's monitoring system and Atlassian Bamboo's API.
It returns the count of how many elastic agents are currently running, consuming an available license.

My company already uses PRTG to monitor systems and processes and we have the SSH Script Advanced Sensor available to us, so I decided to
leverage this system to keep a history of our Elastic Agent Usage

I am also fairly new to Go. There probably are better or more efficient ways to do this,
but this does what I needed it to do so I figured I would share it out there to possibly help others with a similar desire. If you see 
something that I should update to better follow Go standards, please reach out!

## Access Needed

 * Access to create a new sensor in PRTG for your Bamboo sensor (or anywhere that has access to your Bamboo url really)
 * Access to copy the executable and settings.json to your PRTG scriptsxml directory
 * Admin access to your Bamboo server to hit the API

## Configuration

The program expects a 'settings.json' file to be located in the same directory as the executable to provide the following values

 * Username - Username of user with admin access to Bamboo 
 * Password - Password to use the API
 * Url - The entire API url to the remote agents
 * SkipSslVerify - Boolean to ignore SSL cert warnings about the url in the case that it has a self signed certificate, but you really trust it.

It uses Basic Authentication to authenticate with your Bamboo server url.

Example file contents:

```json
{
    "username": "nburglin",
    "password": "somepassword",
    "url": "https://bamboo.mycompany.com:8443/rest/api/latest/agent/remote?online&os_authType=basic"
    "skipsslverify": true
}
```

There is a samle settings.json in this repo that you can overwrite to use. Special characters like '!' work fine with basic auth

##Create Executable

It's expected that you have GO installed locally. You can execute the command below to retrieve the source

```bash
go get github.com/nburglin/BambooElasticAgentCount
```

This will put the executable in your $GOPATH/bin, however you must first make sure you build an executable that will be able to run
on your remote server. For instance, if you are going to be running on a linux host with AMD, your build command will look like this:

```bash
env GOOS=linux GOARCH=amd64 go install github.com/nburglin/BambooElasticAgentCount
```

In the example above, the executable you need will be found in $GOPATH/bin/linux_amd64/BambooElasticAgentCount

Now you simply copy this executable to your remote server to PRTG's scriptsxml directory, along with 
a settings.json mentioned above in the Configuration section. Once this is done, you can create the new SSH Script
Advanced Sensor in PRTG and point it to this executable.


## Flow

This serves a pretty specific purpose, so the flow is fairly simple.

1. Read in config file
2. Set up a Go http client
3. Perform a GET against URL
4. Print out the count of agents with the type of "Elastic" inside of some HTML tags that PRTG expects. The channel name is just hardcoded to "Elastic Agent Count"

## References

More info on PRTG's SSH Script Advanced Sensor found here:
 - [PRTG SSH Script Advanced Sensory](https://blog.paessler.com/prtg-ssh-script-advanced-sensor)

Bamboo API Info (tested as of 6.6.1)
 - [Bamboo API Info](https://docs.atlassian.com/atlassian-bamboo/REST/latest)
