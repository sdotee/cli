# see-cli

A command-line client for the [S.EE](https://s.ee) content sharing platform, enabling users to create and manage short URLs and text snippets and other content efficiently.

## Install

### macOS (Homebrew)

```bash
brew tap sdotee/homebrew-tap
brew install see
```

### Debian / Ubuntu / RHEL (xTom repository)

[xTom](https://xtom.com/) publishes `see` in its [APT and DNF repository](https://github.com/xtomcom/repo) for `amd64` and `arm64`, so updates come through `apt upgrade` / `dnf upgrade`. It supports Debian 12 and 13, Ubuntu 22.04, 24.04 and 26.04 LTS, and RHEL, Rocky Linux and AlmaLinux 9 and 10.

```bash
# Add the repository (DEB822 sources or a .repo file, as appropriate)
curl -sL ba.sh/repo | sudo bash

sudo apt install see   # Debian / Ubuntu
sudo dnf install see   # RHEL / Rocky / AlmaLinux
```

Manual setup:

```bash
# Debian / Ubuntu
sudo apt install -y ca-certificates curl gnupg
curl -fsSL https://repo.xtom.com/xtom.key | sudo bash -c 'gpg --dearmor > /usr/share/keyrings/xtom.gpg'
echo "deb [signed-by=/usr/share/keyrings/xtom.gpg] https://repo.xtom.com/deb stable main" | sudo tee /etc/apt/sources.list.d/xtom.list
sudo apt update && sudo apt install see

# RHEL / Rocky / AlmaLinux
sudo curl -o /etc/yum.repos.d/xtom.repo https://repo.xtom.com/rpm/xtom.repo
sudo dnf install see
```

### Linux packages

Each [release](https://github.com/sdotee/cli/releases) has `.deb`, `.rpm`, `.apk` and `.pkg.tar.zst` packages for `amd64` and `arm64`:

```bash
sudo dpkg -i see-cli_<version>_linux_<arch>.deb                    # Debian / Ubuntu
sudo rpm -i see-cli_<version>_linux_<arch>.rpm                     # RHEL / Fedora
sudo apk add --allow-untrusted see-cli_<version>_linux_<arch>.apk  # Alpine
sudo pacman -U see-cli_<version>_linux_<arch>.pkg.tar.zst          # Arch Linux
```

### Arch Linux (AUR)

```bash
yay -S see-cli       # build from source
yay -S see-cli-bin   # pre-built binary
```

### Pre-built binaries

Releases also include `see_<OS>_<arch>.tar.gz` for macOS, Linux and FreeBSD, and `see_Windows_<arch>.zip` for Windows, on `x86_64` and `arm64`.

### From source

Requires Go 1.26+:

```bash
git clone https://github.com/sdotee/cli.git
cd cli
make build   # builds ./see
```

## Configuration

| Flag         | Environment Variable | Description           |
| ------------ | -------------------- | --------------------- |
| `--api-key`  | `SEE_API_KEY`        | API key (Required)    |
| `--base-url` | `SEE_BASE_URL`       | API base URL          |
| `--timeout`  | `SEE_TIMEOUT`        | Request timeout (flag: duration such as `30s`; env: seconds) |
| `--json`     |                      | Output in JSON format |

Run `see help <command>` for the complete flag reference. Add `--json` to any command to receive machine-readable output.

## Commands

### Domains & Tags

List available domains and tags:

```bash
see domains
see tags
```

### Short URLs

Manage short links. `--domain` defaults to `s.ee`; `--expire-at` takes a Unix timestamp in seconds.

**Create**

```bash
see shorturl create <url> [flags]

# Flags:
# --slug, --domain, --title, --password, --expire-at, --tag-ids, --expiration-redirect-url
```

**Update**

```bash
see shorturl update <slug> --target-url <url> [--domain s.ee] [--title <title>]
```

**Delete**

```bash
see shorturl delete <slug> [--domain s.ee]
```

**History and statistics**

```bash
see shorturl history --page 1
see shorturl stats <slug> --domain s.ee --period daily
```

The statistics period can be `daily`, `monthly`, or `totally`.

### Text

Manage text snippets. Reads from stdin by default or `--file`. `--type` can be `plain_text`, `source_code` or `markdown`.

**Create**

```bash
echo "hello" | see text create [flags]

# Flags:
# --file, --type, --slug, --domain, --title, --password, --expire-at, --tag-ids
```

**Update**

```bash
see text update <slug> --file <path> [--domain s.ee] [--title <title>]
# or pipe the new content via stdin
```

**Delete**

```bash
see text delete <slug> [--domain s.ee]
```

**Domains and history**

```bash
see text domains
see text history --page 1
```

### Files

Upload and manage files.

**Domains**

List available domains for file uploads:

```bash
see file domains
```

**Upload**

```bash
see file upload [files...] [flags]
# OR
cat image.png | see file upload --name image.png

# Flags:
# --file, -f: Path to file (optional if passed as argument)
# --name, -n: Filename (required if using stdin)
# --domain: Domain for the file URL
# --slug: Custom slug for the file URL
# --private: Upload as a private file
```

Regular files up to 100 MB use multipart upload. For seekable files larger than 100 MB, the CLI automatically uses the SDK's resumable TUS upload flow, supporting files up to 5 GB. Stdin has no discoverable size and uses regular upload.

**History**

List uploaded file history (30 files per page):

```bash
see file history [flags]

# Flags:
# --page, -p: Page number (default 1)
```

**Download URL**

Get a temporary download URL for a private file:

```bash
see file download-url <file-id>
```

The returned URL is temporary and normally remains valid for about one hour.

**Delete**

```bash
see file delete <delete_keys...>
```

### Bio Pages

```bash
see bio create "My Bio" \
	--description "About me" \
	--mastodon-url https://c.im/@s.ee \
	--link Blog=https://example.com

see bio update <id> --title "Updated Bio" --link Blog=https://example.com
see bio history --page 1
see bio delete <id>
```

`--link title=url` can be repeated. `update` requires `--title`, and links supplied during an update replace the existing custom links. `create` also accepts `--domain`, `--slug` and `--rss-url`.

### QR Codes

```bash
see qrcode create https://example.com --title "Example"   # --title is required
see qrcode history --page 1
see qrcode delete <slug> --domain s.ee
```

Creation returns PNG, SVG, and PDF URLs generated by S.EE.

### Account

Show API, content, file, and storage usage:

```bash
see usage
```

Validate the configured API key or another token:

```bash
see token check
see token check <token>
```

### Version

```bash
see version
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
