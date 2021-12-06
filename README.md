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

The build will be made available under **dist**.

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
command=/var/www/path/to/recipya/dist/recipya serve
directory=/var/www/path/to/recipya/dist
autorestart=true
autostart=true
stdout_logfile=/var/www/path/to/recipya/logs/supervisord.log

$ sudo supervisorctl
> status
> update
> status
```

## Recipes Database

A folder of recipes is required for this application to work when running locally because
the database will index them.

The folder can be placed anywhere. Each recipe is a JSON-formatted file that follows the [recipe schema](https://schema.org/Recipe) standard.
Not all fields are currently supported. Refer to the [Feedback](#feedback) if you require a field from the schema.

If you use Nextcloud Cookbook, then all is needed is to point to the folder where Cookbook stores the recipes.

### No Recipes Right Now?

No problem! Download this [sample](https://cloud.musicavis.ca/index.php/s/NNge4Dp7sHXrPsW).

## Docker

Run the following Docker command to run the project as a daemon in a container named "recipya".

```bash
$ docker run -d \
   --name recipya \
   --restart=unless-stopped \
   -v /path/to/your/json-recipes/repository:/recipes \
   -p 3001:3001 \
   -e RH_INDEXINTERVAL=2d \
   -e RH_WAIT 10 \
   reap99/recipya:latest
```

Finally, give the API a try:

```bash
$ curl "localhost:3001/api/v1/recipes/search?ingredients=avocado&num=1"
```

## Docker-Compose

Download the [docker-compose.yaml](https://github.com/reaper47/recipya/blob/main/docker-compose.yaml) file along with its [configuration](https://github.com/reaper47/recipya/blob/main/docker-compose.yaml) file.

```bash
$ curl -LJO https://raw.githubusercontent.com/reaper47/recipya/main/docker-compose.yaml -LJO https://raw.githubusercontent.com/reaper47/recipya/main/.env
```

Modify the configuration variables in the `.env` file if needed.

```bash
$ nano .env
```

Then, run the container:

```bash
$ docker-compose up -d
```

## Environment Variables

The following environment variables can be set if you want to use a value different from the default:

### RH_INDEXINTERVAL

The interval at which the recipes database will be indexed. A value of -1 will disable the automatic indexing [Note: To be implemented].

The syntax is [int][unit], e.g. 10m, 1h, 1d.

Valid units are:

- m: minutes
- h: hours
- d: days
- M: months
- w: weeks
- y: years

Default: 1d

### RH_WAIT

The duration in seconds for which the server gracefully waits for existing connections to finish.

Default: 10

# Deployment

Here is a sample Caddy blob to expose the container to the outside world:

```bash
recipes.your-domain.name {
    encode zstd gzip
    reverse_proxy localhost:3001

    header Access-Control-Allow-Method "GET, OPTIONS"
    header Access-Control-Allow-Headers "*"
    header Access-Control-Allow-Origin "*"

    log {
        output file /var/log/caddy/recipes.your-domain.name.access.log
    }
}
```

**Note:** The access control headers might move to the server in the future.

Then, reload the server:

```bash
$ caddy reload
```

No Nginx configuration is given because Caddy is easy, simple, has automatic HTTPS, generates and renews certbot certificates automatically, and works like a charm.

# Run Locally

1. Clone the project:

```bash
$ git clone https://github.com/reaper47/recipya.git
```

2. Move inside the directory:

```bash
$ cd recipya
```

3. Install the dependencies:

```bash
$ go mod download
```

4. Build the program:

```bash
$ go build -o dist
```

5. Modify the configuration file in /dist

6. You are ready to go!

# Usage/Examples

To display the help text:

```bash
$ ./recipya
```

To index the database:

```bash
$ ./recipya help index
$ ./recipya index
```

To search for recipes from a list of ingredients:

```bash
$ ./recipya help search
$ ./recipya search avocado,garlic,ketchup
```

To start the web server that exposes the REST API:

```bash
$ ./recipya help serve
$ ./recipya serve
$ curl localhost:3001/api/v1/recipes/search?ingredients=avocado,garlic,ketchup&num=3&mode=2
```

# Running Tests

To run tests, run the following command:

```bash
$ go test ./...
```

# API Reference

## Search recipes

```http
  GET /api/v1/search
```

| Parameter     | Type     | Description                                                                                                                                                                                                          |
| :------------ | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ingredients` | `string` | **Required**. A comma-separated list of ingredients                                                                                                                                                                  |
| `num`         | `int`    | The maximum number of recipes to fetch.<br>Default: `10`.                                                                                                                                                            |
| `mode`        | `int`    | The search mode to employ.<br>Mode `1`: Minimize the number of missing ingredients in order to buy less at the grocery store.<br>Mode `2`: Maximize the number of ingredients taken from the fridge.<br>Default `2`. |

# Tech Stack

**Client:** Vue.js

**Server:** Go

**Deployment:** Docker

# Feedback

If you have any feedback, please reach out to us at macpoule@gmail.com or open an issue on GitHub.
