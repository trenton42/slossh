# Slossh - your friendly, extensible ssh sentinel

Slossh is a simple ssh server that never lets you in. It's purpose is to sit on your network, collecting information about bots from bad actors who are constantly attempting to gain access to your systems. It also allow you to make real time decisions, blocking access to your network to IPs that are currently engaging in attacks. It can also be used to research username / password combinations that are being used by these bots to ensure that your password policies are effective.

## Placement

Ideally, the program should be run on a device in your ip space, but separate from the rest of your network by policy. This will allow blocking of active nefarious IPs while slossh continues to gather data. It should listen on port 22, but run as an unprivileged user.

## Recorders

Internally, Slossh has recorders that take the data that is collected and make it useful. They might store the data in a file, or send it to a remote endpoint for immediate action. There are currently two built in recorders (`file` and `http post`), but more are planned, and it is fairly simple to add a new recorder.

### File recorder

Store json results locally in a single file.

### HTTP recorder

Send results to a remote http server. This recorder sends a `POST` request to the specified URL with a json encoded body for each login attempt.

# SSH Session JSON format

Each session may contain zero or more login attempts. The session structure is sent to each configured recorder once it is closed. If public key authentication is attempted, the public key used will be recorded.

```json
{
    "SessionID": "cf3750f6a7bf951fc2aa5a3c05bd8aa7050b94c7f7e5d7d09afa18bf20b7e2d2",
    "IP": "127.0.0.1",
    "ClientVersion": "SSH-2.0-OpenSSH_8.1",
    "Attempts": [
        {
            "Username": "admin",
            "Key": {
                "Key": "ssh-rsa AAA..... # full ssh public key here",
                "Fingerprint": "SHA256:s6iMqZs5Uh8x530Sjlwqes9m/w1UykbK0x29pfupPSo",
                "Type": "ssh-rsa"
            },
            "Password": ""
        },
        {
            "Username": "admin",
            "Key": null,
            "Password": "P@SSw0rd"
        }
    ],
    "Start": "2020-11-20T22:23:26.820267-05:00",
    "Finish": "2020-11-20T22:23:32.139888-05:00"
}
```

## Command line options

```
    --file-path string       Path to json file to store results
    --http-url string        URL to send post requests to
-p, --port int               Port to listen on (default 2022)
-r, --recorder stringArray   recorder to use (can be specified multiple times). Available recorders: file, http
  ```

## License

MIT