# watchlog

Tail and filter structured JSON logs from multiple sources with real-time field extraction and colorized output.

---

## Installation

```bash
go install github.com/yourusername/watchlog@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/watchlog.git
cd watchlog
go build -o watchlog .
```

---

## Usage

Tail a single log file and filter by field value:

```bash
watchlog --file /var/log/app.json --filter level=error
```

Watch multiple sources and extract specific fields:

```bash
watchlog --file /var/log/app.json --file /var/log/worker.json \
  --fields time,level,message \
  --filter level=warn
```

Pipe from stdin:

```bash
kubectl logs -f my-pod | watchlog --fields time,level,message --filter level=error
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Log file to tail (repeatable) |
| `--filter` | Filter by field value (e.g. `level=error`) |
| `--fields` | Comma-separated list of fields to display |
| `--no-color` | Disable colorized output |

---

## License

MIT © [yourusername](https://github.com/yourusername)