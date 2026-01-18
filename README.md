
```markdown
age-wrapper — minimal age CLI wrapper

Overview
--------
age-wrapper is a small, user-friendly CLI wrapper around the `age` tool (https://filippo.io/age/).
It provides these commands:

- `agewrap keygen <out-file>` — generate an `age` identity (wraps `age-keygen`)
- `agewrap encrypt -r <recipient> [-o out] <infile>` — encrypt for one or more recipients (supports `github:username`)
- `agewrap decrypt -i <identity> [-o out] <infile>` — decrypt using an identity file

It shells out to the system `age` and `age-keygen` binaries (keeps code minimal).

Requirements
------------
- Go 1.20+
- `age` and `age-keygen` installed and in PATH

Install `age`:
- macOS: `brew install age`
- Ubuntu/Debian: `apt install age` (or build from source)
- Official: https://filippo.io/age/

Quickstart
----------
1. Save files from this repo into a folder (module path is github.com/Charan-Hari/age-wrapper).
2. Build:

   go mod tidy
   go build -o agewrap ./...

3. Generate an identity:

   ./agewrap keygen -p mykey.txt

4. Encrypt:

   ./agewrap encrypt -r "age1..." -o secret.txt.age secret.txt

   Or use GitHub keys:
   ./agewrap encrypt -r github:alice -o secret.txt.age secret.txt

5. Decrypt:

   ./agewrap decrypt -i mykey.txt -o secret.txt secret.txt.age

Notes
-----
- `github:username` fetches https://github.com/username.keys and uses returned public keys as recipients.
- If your `age` binary does not support SSH recipients, you may need to update `age` or convert the key.

Packaging
---------
Run:

chmod +x package.sh
./package.sh

A `age-wrapper.zip` will be created.

```
