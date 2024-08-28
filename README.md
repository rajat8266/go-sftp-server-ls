## SFTP Server with Local Disk as File Storage
This project implements an SFTP server that uses local disk as the backend storage system. It provides a secure, scalable, and configurable SFTP service that can be used for file transfers.

### Sftp Backend
**Overview**: The SFTP server is built using the Go programming language, with support for secure file transfer over SSH. It handles file operations such as reading, writing, listing, and deleting files.

**Key Components**:
- **Connection Handling**: Manages incoming SSH connections using `golang.org/x/crypto`, allowing for secure file transfer sessions.
- **SFTP Handlers**: Custom handlers for SFTP operations are implemented, using `github.com/pkg/sftp`.
- **Session Management**: Supports multiple concurrent SFTP sessions, with each session linked to specific storage location based on user authentication.

### Storage
**Overview**: The project uses disk space or folder path as the primary data store, where all files uploaded via SFTP are stored and managed. Program will automatically create folders for the users defined in config.json file if not found already.

**Key Components**:
- **User-Specific Folder**: Each user is associated with a specific folder, ensuring data isolation and security.


### Configurations
**Overview**: The project is highly configurable, allowing administrators to define key settings via a JSON configuration file and environment variables.

**Key Components**:
- **Root Configuration**: Loads and parses a JSON configuration file that specifies settings such as the address and port for the SFTP server and user accounts.
- **SSH Server Configuration**: Configures the SSH server with the necessary settings, including private keys and authorized public keys for user authentication.


### Run
```./go-sftp-server --config-path path_to_config.json```
Generate a new config using `config/config.md`.

