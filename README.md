# RemoteLog (WIP)

RemoteLog is a command-line tool written in Go that connects to a remote server and tails log files. It provides functionality to monitor and filter log entries based on specified criteria. RemoteLog supports connecting using either a private key or a Kerberos ticket.

## Installation

To install RemoteLog, follow these steps:

1. Make sure you have Go installed on your system. If not, you can download it from the official Go website: [https://golang.org/](https://golang.org/)

2. Clone the RemoteLog repository to your local machine:

```bash
git clone https://github.com/a-dakani/remoteLog.git
```

3. Navigate to the RemoteLog directory:

```bash
cd remoteLog
```

4. Build the RemoteLog binary:

```bash
make build
```

5- The binary and the configuration file will be available in the `/build` directory.

## Usage

RemoteLog can be used to monitor log files on a remote server. It provides various command-line options to configure the connection and filtering options. Here are the available command-line flags:

| Flag     | Description                                                              |
| -------- | ------------------------------------------------------------------------ |
| `-srv`   | Specify a predefined service name from the `config.services.yaml` file.  |
| `-fs`    | Provide a list of log file paths to tail, separated by commas.            |
| `-h`     | Set the host to connect to.                                              |
| `-u`     | Set the user to connect to the host.                                      |
| `-p`     | Set the port to connect to the host. Defaults to `22`.                    |
| `-pk`    | Set the path to the private key for authentication.                       |
| `-krb5`  | Set the location of the `krb5.conf` file for Kerberos authentication.     |
| `-f` (WIP) | Set filter words for log file entries. Multiple words should be separated by commas. |

### Examples

Here are a few examples to demonstrate how RemoteLog can be used:

1. Tailing log files using predefined service:

```bash
./build/remoteLog -srv myService
```

2. Tailing log files using command-line arguments:

```bash
./build/remoteLog -fs=/var/log/app.log,/var/log/error.log -h=192.168.1.1 -u=admin -p=22 -pk=/path/to/private/key -f=ERROR,WARN
```

## Configuration

RemoteLog relies on configuration files to define services and their respective settings. The configuration files should be placed in the same directory as the RemoteLog binary and have the following formats:

- `/build/config.yaml`: Contains general configurations for RemoteLog.
- `/build/config.services.yaml`: Defines predefined services and their connection details.

Please make sure to configure these files correctly before running RemoteLog.
Here's an example of what the file should look like:

```yaml

services:
  - name: spyName
    host: spyHost.live.com
    user: spyUser
    port: 22
    #private_key_path: ~/.ssh/id_rsa   
    krb5_conf_path: /etc/kerb5.conf
    files:
      - alias: spy-1
        path: /var/log/spyPath.log
      - alias: spy-2
        path: /var/log/spyPath.log
```

- choose your Auth-Method for each Service between `private_key_path` and `krb5_conf_path`
- If both are provided, the private_key_path will be used first.
- If it fails, the krb5_conf_path will be used.


## Known Issues:
- Error logging for Kerberos tickets is currently buggy and may not provide accurate error messages or handle certain scenarios correctly. This can affect the authentication process when connecting to a remote server using a Kerberos ticket.
