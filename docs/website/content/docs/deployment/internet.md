---
title: Web Hosting
weight: 2
next: /docs/development
---

## Caddy

[Caddy](https://caddyserver.com/) is a lightweight, extensible open-source web server that
automatically obtains and renews TLS certificates for all your sites.

The following block in the Caddyfile is the minimum required for hosting Recipya over the network.
If you are on Linux, please ensure that Recipya is running as a [service](/guide/docs/deployment/local-network/#linux).

```text
$ sudo cat /etc/caddy/Caddyfile

domain.com {
	encode zstd gzip
	reverse_proxy localhost:PORT

	header /static/* Cache-Control "public, max-age=2678400, must-revalidate"

	log {
		output file /var/log/caddy/domain.com.access.log
	}
}
```
