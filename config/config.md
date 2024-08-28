## Configuration guide for go-sftp-server
This document provides an overview of the configuration settings required for the fo-sftp-server application. These settings are defined in a JSON configuration file and are loaded at runtime.

### Configuration Structure
The main configuration file is parsed into a Config struct, which holds various settings for the application, including paths to SSH keys, SFTP accounts.

Example Configuration File:

```json
{
  "SftpAccounts": [
    {
      "Username": "user1",
      "Password": "password1",
    },
    {
      "Username": "user2",
      "Password": "password2",
    }
  ],
  "Address": "127.0.0.1",
  "Port" : "2024",
  "SSHKeyPath": "**required**",
  "BasePath" : "**required**",
  "SftpAuthorizedKeysFile": "",
}
```

### Config Fields
- **SSHKeyPath**:
Path to the SSH private key used by the SFTP server for secure connections.
Generate using `ssh-keygen -b 2048 -t rsa -f filepath.txt`

- **BasePath**:
Path to disk where data will be stored. For "user1", Files will be stored at location BasePath/Username. Program should have read/write access to BasePath.

- **SftpAuthorizedKeysFile**:
Path to the file containing authorized SSH public keys for users who are allowed to connect. If this field is not set, only password authentication will be used.

- **SftpAccounts**:
List of SFTP accounts, where each account contains the following fields:

  - **Username**:
  The username for the SFTP account.

  - **Password**:
  The password for the SFTP account. This is used for password-based authentication.



