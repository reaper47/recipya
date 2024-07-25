---
title: Hébergement web
weight: 2
next: /docs/development
---

## Apache

Le [serveur HTTP Apache](https://httpd.apache.org) est un logiciel de serveur web open source qui set du contenu web sur l'internet.
Il est largement utilisé en raison de sa fiabilité, de sa flexibilité et de son extensibilité dans la prise en charge de divers
technologies du Web.

Le bloc suivant montre le minimum requis dans le fichier de configuration d'Apache pour héberger Recipya sur le réseau.

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

[Caddy](https://caddyserver.com/) est un serveur web open source léger et extensible qui obtient et renouvelle automatiquement les certificats TLS pour tous vos sites.

Le bloc suivant montre le minimum requis dans le fichier de configuration de Caddy, nommé Caddyfile, pour héberger Recipya sur le réseau.
Si vous utilisez Linux, assurez-vous que Recipya s'exécute en tant que [service](/guide/docs/deployment/local-network/#linux).

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

## Nginx

[Nginx](https://en.wikipedia.org/wiki/Nginx) est un serveur web puissant qui peut également être utilisé comme proxy inverse, équilibreur de charge, proxy de messagerie et cache HTTP.
Il est largement utilisé pour ses hautes performances, son efficacité dans la gestion des connexions simultanées et sa faible consommation de ressources.

Le bloc suivant montre le minimum requis dans le fichier de configuration de Nginx pour héberger Recipya sur le réseau.

```text
server {   
    listen 80;
    server_name domain.com;

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript;
     
    location / {
        proxy_pass http://127.0.0.1:<port>/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location ~* /static/ {
        add_header Cache-Control "public, max-age=2678400, must-revalidate";
    }
    
    location /ws {
        proxy_pass http://127.0.0.1:8125;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```