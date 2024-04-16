---
title: Web Hosting
weight: 2
next: /docs/development
---

## Apache

The [Apache HTTP Server](https://httpd.apache.org) is an open-source web server software that serves web content 
over the internet. It's widely used due to its reliability, flexibility, and extensibility in supporting various 
web technologies.

The following block in the Apache configuration file is the minimum required for hosting Recipya over the network.

```text
<IfModule mod_ssl.c>
<VirtualHost *:443>
    ServerAdmin [email]
    ServerName [subdomain.domain.com]

    ProxyPass / http://127.0.0.1:<port>/
    ProxyPassReverse / http://127.0.0.1:<port>/

    RewriteEngine on
    RewriteCond %{HTTP:UPGRADE} ^WebSocket$ [NC]
    RewriteCond %{HTTP:CONNECTION} Upgrade$ [NC]
    RewriteRule .* ws://127.0.0.1:<port>%{REQUEST_URI} [P]

    ErrorLog ${APACHE_LOG_DIR}/[log file]
    CustomLog ${APACHE_LOG_DIR}/[log file]

    SSLCertificateFile [letsencrypt file]
    SSLCertificateKeyFile [letsencrypt file]
    Include [letsencrypt ssl file]
</VirtualHost>
</IfModule>
```

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
