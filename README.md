# subfalcon

subfalcon is a subdomain enumeration tool that allows you to discover and monitor subdomains for a given list of domains or a single domain. It fetches subdomains from various sources, checks for potential subdomain takeover vulnerabilities, saves findings to a SQLite database, and can notify updates via Discord.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)
- [ToDo](#todo)

## Features

- Subdomain enumeration from multiple sources:
    1. crt.sh
    2. hackertarget
    3. anubis
    4. Alienvault
    5. rapiddns
    6. urlscan.io
- Subdomain takeover scanning
    - Currently supports Azure services (cloudapp.net, azurewebsites.net, cloudapp.azure.com)
    - Colored terminal output for scan results
- SQLite database to store discovered subdomains
- Enhanced Discord integration
    - Separate notifications for new subdomains and takeover findings
    - Formatted messages with emojis and proper formatting
    - Automatic file attachments for large result sets
    - Rate limit handling
- Easy-to-use command-line interface
- Option to process a single domain with the `-d` flag

## Installation
You can install subfalcon using the following command: 
```bash
go install github.com/cyinnove/subfalcon/cmd/subfalcon@latest
```

## Usage

```bash
subfalcon -d example.com -sdt -m -wh "YOUR_DISCORD_WEBHOOK_URL"
```

## Options

- `-l` or `--domain_list`: Specify a file containing a list of domains
- `-m` or `--monitor`: Monitor subdomains and send updates to Discord
- `-wh` or `--webhook`: Specify the Discord webhook URL
- `-d` or `--domain`: Specify a single domain for processing
- `-sdt`: Enable subdomain takeover scanning

## Examples

- Basic usage with subdomain takeover scanning:
  ```bash
  subfalcon -d example.com -sdt
  ```

- Monitor a single domain with takeover scanning and Discord notifications:
  ```bash
  subfalcon -d example.com -sdt -m -wh "YOUR_DISCORD_WEBHOOK_URL"
  ```

- Monitor multiple domains with all features:
  ```bash
  subfalcon -l domains.txt -sdt -m -wh "YOUR_DISCORD_WEBHOOK_URL"
  ```

## Contributing

Feel free to contribute by opening issues or submitting pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

## Disclaimer

Use this tool responsibly and only on systems you have permission to scan. The authors are not responsible for any misuse or damage caused by this tool.

## ToDo

- [x] Add subdomain takeover scanning
- [x] Improve Discord notifications with better formatting
- [x] Add file attachment support for large result sets
- [ ] Add support for more takeover vulnerability patterns
- [ ] Add monitoring using Telegram
- [ ] Add more subdomain enumeration sources
- [ ] Add flags to customize monitoring time intervals
- [ ] Add concurrency for faster subdomain enumeration
- [ ] Add proxy support for requests
- [ ] Add custom output formats (JSON, CSV)
- [ ] Add vulnerability severity levels
- [ ] Add support for custom takeover patterns
- [ ] Improve error handling and logging system

> If you enjoy what we do, please support us:
> Buy Me Ko-fi! https://ko-fi.com/h0tak88r
