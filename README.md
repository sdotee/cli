# see-cli

A command-line client for the [S.EE](https://s.ee) content sharing platform, enabling users to create and manage short URLs and text snippets and other content efficiently.

## Install

```bash
go install .
```

On Debian/Ubuntu:
```bash
dpkg -i see-cli_<version>_<arch>.deb
```

On RHEL/Fedora:
```bash
rpm -i see-cli-<version>-1.<arch>.rpm
```

## Configuration

| Flag         | Environment Variable | Description           |
| ------------ | -------------------- | --------------------- |
| `--api-key`  | `SEE_API_KEY`        | API key (Required)    |
| `--base-url` | `SEE_BASE_URL`       | API base URL          |
| `--timeout`  | `SEE_TIMEOUT`        | Request timeout       |
| `--json`     |                      | Output in JSON format |

## Commands

### Domains & Tags

List available domains and tags:

```bash
see domains
see tags
```

### Short URLs

Manage short links.

**Create**

```bash
see shorturl create <url> [flags]

# Flags:
# --slug, --domain, --title, --password, --expire-at, --tag-ids, --expiration-redirect-url
```

**Update**

```bash
see shorturl update <slug> --target-url <url> [flags]
```

**Delete**

```bash
see shorturl delete <slug>
```

### Text

Manage text snippets. Reads from stdin by default or `--file`.

**Create**

```bash
echo "hello" | see text create [flags]

# Flags:
# --file, --type, --slug, --domain, --title, --password, --expire-at, --tag-ids
```

**Update**

```bash
see text update <slug> [flags]
```

**Delete**

```bash
see text delete <slug>
```

### File Upload

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
# --is-private: Whether this file should be private (0 = public, 1 = private)
```

**History**

List uploaded file history (30 files per page):

```bash
see file history [flags]

# Flags:
# --page, -p: Page number (default 1)
```

**Delete**

```bash
see file delete <delete_keys...>
```

### Version

```bash
see version
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
