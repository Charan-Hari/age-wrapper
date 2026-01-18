```markdown
# age-wrapper — minimal age CLI wrapper

![Go version](https://img.shields.io/badge/go-1.20+-00ADD8)

Overview
--------
`age-wrapper` (binary: `agewrap`) is a tiny, user-friendly CLI wrapper around the `age` tool (https://filippo.io/age/). It provides a minimal local workflow:

- `agewrap keygen <out-file>` — generate an `age` identity (wraps `age-keygen`)
- `agewrap encrypt -r <recipient> [-o out] <infile>` — encrypt for one or more recipients (supports `github:username`)
- `agewrap decrypt -i <identity> [-o out] <infile>` — decrypt using an identity file

This project intentionally shells out to the system `age` and `age-keygen` binaries to keep the code small and portable.

Requirements
------------
- Go 1.20+ (to build)
- `age` and `age-keygen` installed and on your `PATH`

Install `age`:
- macOS (Homebrew):
  ```bash
  brew install age
  ```
- Debian / Ubuntu:
  ```bash
  sudo apt install age
  ```
  or build from source per the official instructions.

Official project: https://filippo.io/age/

Repository
----------
This project is prepared for the module path:

```
github.com/Charan-Hari/age-wrapper
```

Installation (build locally)
----------------------------
Clone or copy the project files into a directory and run:

```bash
# make sure go is installed and in PATH
go mod tidy
go build -o agewrap ./...
```

The produced binary will be `./agewrap`.

Quickstart
----------
1. Generate an identity (secret key file):

```bash
./agewrap keygen -p mykey.txt
# -p/--show-public prints the public key comment line in the generated file
```

The file `mykey.txt` will contain something like:
```
# public key: age1...
AGE-SECRET-KEY-1...
```

2. Encrypt a file for a recipient (age recipient string):

```bash
./agewrap encrypt -r "age1..." -o secret.txt.age secret.txt
```

3. Encrypt using GitHub public keys (convenience: `github:username`):

```bash
./agewrap encrypt -r github:alice -o secret.txt.age secret.txt
```

This fetches `https://github.com/alice.keys` and uses any returned public keys as recipients.

4. Decrypt using your identity file:

```bash
./agewrap decrypt -i mykey.txt -o secret.txt secret.txt.age
```

Usage examples
--------------
Encrypt for multiple recipients (mix of GitHub and direct recipients):

```bash
./agewrap encrypt \
  -r github:alice \
  -r "age1xxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  -o multi.age \
  somefile.txt
```

Decrypt an armored file (if `age` supports armor) — same decrypt command works:

```bash
./agewrap decrypt -i mykey.txt -o plaintext.txt secret.txt.age
```

Packaging
---------
A helper script `package.sh` is provided to create a distributable zip:

```bash
chmod +x package.sh
./package.sh
# produces age-wrapper.zip
```

Notes & troubleshooting
-----------------------
- `github:username` support:
  - The tool fetches the public keys at `https://github.com/<username>.keys` and passes them to `age` as recipients.
  - Newer `age` versions accept OpenSSH public keys (e.g. `ssh-ed25519 ...`) as recipients. If your `age` binary rejects an SSH key recipient, either:
    - Update `age` to a newer release that supports SSH recipients, or
    - Convert the SSH public key to an X25519 `age` recipient (see age docs), or
    - Use `age` public recipients directly (`age1...`).

- If you see `required executable "age" not found in PATH` or `age-keygen not found`, install `age` and ensure the binaries are accessible from your shell.

- The wrapper prints the exact `age` command it runs — you can copy that command to debug or run manually.

Security notes
--------------
- This wrapper does not implement cryptography itself: it calls the well-tested `age` binary.
- Keep your AGE secret keys (`AGE-SECRET-KEY-1...`) securely stored (file permissions, encrypted vaults, etc.).
- Do not commit secret key files to version control.

Contributing
------------
Contributions are welcome. Suggested improvements:
- Add automatic ssh-ed25519 → age X25519 conversion for broader compatibility.
- Add tests and CI (GitHub Actions) for build and linting.
- Add an install/release pipeline to produce platform-specific binaries.


Acknowledgements
----------------
- `age` by Filippo Valsorda — https://filippo.io/age/
- Cobra CLI library — https://github.com/spf13/cobra
```
