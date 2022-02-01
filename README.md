# Recipya

Recipya is an application whose goal is to search for what you can cook with the ingredients in your fridge.
In other words, it helps you know what you can cook with what you have when you are out of ideas.

It works seamlessly with recipes in your [Nextcloud Cookbook](https://apps.nextcloud.com/apps/cookbook).

The project consists of a backend and a frontend.
The backend is a REST API. The frontend, found under the /web folder, is a simple app where the user can use the search function.

# Features

- Search for recipes based on what you have in your fridge
- Cross-platform solution
- Can be self-hosted

# Installation

## Database

Recipya uses PostgreSQL to store data.

Install (Debian)
```bash
$ sudo apt-get install postgresql postgresql-contrib
```

Enable PostgreSQL on start:
```bash
$ sudo systemctl enable postgresql
```

Create the database:
```bash
$ sudo su - postgres
$ psql
$ CREATE USER recipya WITH password 'elephants';
$ CREATE DATABASE recipya OWNER recipya;
```

## Recipya

Clone the repository:
```bash
$ git clone https://github.com/reaper47/recipya.git
```

Build/update the program:
```bash
$ sudo sh update.sh
```

The build will be made available under **bin**.

## Self-host

Caddy server:
```bash
$ sudo nano /etc/caddy/Caddyfile

...
domain {
	encode zstd gzip

	header /static/* Cache-Control "public, max-age=2678400, must-revalidate"
	
	log {
		output file /var/www/path/to/recipya/logs/caddy-access.log
		format single_field common_log
	}

	reverse_proxy http://localhost:8080
}
...

$ sudo mkdir /var/www/path/to/recipya/logs
```

Supervisor to start Recipya as a daemon:
```bash
$ sudo nano /etc/supervisor/conf.d/recipya.conf

[program:recipya]
command=/var/www/path/to/recipya/bin/recipya serve
directory=/var/www/path/to/recipya/bin
autorestart=true
autostart=true
stdout_logfile=/var/www/path/to/recipya/logs/supervisord.log

$ sudo supervisorctl
> status
> update
> status
```

# Running Tests

To run tests, run the following command:

```bash
$ make test
```

# Feedback

If you have any feedback, please reach out to us at macpoule@gmail.com or open an issue on GitHub.
