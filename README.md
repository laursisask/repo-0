# Keeper Security's Linux Keyring Utility

This utility interacts with the native Linux APIs to store and retrieve secrets from the Keyring using [Secret Service](https://specifications.freedesktop.org/secret-service/latest/).

While initially developed to help Keeper secure KSM configs, this utility can be used by any integration, plugin, or code base, to store and retrieve credentials, secrets, and passwords in the Linux Keyring simply and natively.

## Setup 

Download the latest version from the releases page and optionally add it to PATH to get started.

## Usage

The executable supports two commands:

1. `set`
2. `get`

Both commands require an application `name` (i.e. the name of the secret in / to be stored in the Keyring) as the first argument.

### `set`

`set` requires a second argument of the secret to be stored. This can be either a:

1. BASE64 string
2. JSON string
3. Path to an existing JSON file

When the secret is saved to the Keyring it is first encoded into a BASE64 format (if not already a BASE64 string). This standardizes the format for both consistent storage and to make it easier to consume by Keeper integrations and products. 

> If you need a support for a different format, please submit a feature request. We'd be happy to extend this to support other use cases.

### `get`

`get` returns the stored BASE64 encoded config to `stdout` and exits with a `0` exit code. The requesting integration can capture the output for consumption. Any errors encountered retrieving the config will return a `non-zero` exit code and write to `stderr`.

### Example usage

```shell
# Save a secret
lku set APPNAME eyJ1c2VybmFtZSI6ICJnb2xsdW0iLCAicGFzc3dvcmQiOiAiTXlQcmVjaW91cyJ9
# or
lku set APPNAME config.json

# Retrieve a secret
lku get APPNAME
```

## Contributing

Please read and refer to the contribution guide before making your first PR.

For bugs, feature requests, etc., please submit an issue!