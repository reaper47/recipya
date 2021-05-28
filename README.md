
# Recipe Hunter

A brief description of what this project does and who it's for


## Demo

Insert gif or link to demo

  
## Features

- Light/dark mode toggle
- Live previews
- Fullscreen mode
- Cross platform

  
## Installation 

Install my-project with npm

```bash 
  npm install my-project
  cd my-project
```

## Environment Variables

1. RH_WAIT

1. RH_INDEXINTERVAL

docker run -e RH_WAIT=20 -e RH_INDEXINTERVAL=5m -v /path/to/JSON_recipes/folder:/recipes recipe-hunter

```bash
$ docker run -d \
   --name recipe-hunter \
   --restart=unless-stopped \
   --user $(id -u):$(id -g) \
   -v /path/to/json-recipes:/recipes \
   -p 3001:3001 \ 
   -e RH_INDEXINTERVAL=1d \
   recipe-hunter
   reap47/recipe-hunter:latest
```
docker run -d --name recipe-hunter --restart=unless-stopped -v /path/to/json-recipes:/recipes -p 3001:3001 -e RH_INDEXINTERVAL=1d recipe-hunter

## Run Locally

Clone the project

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd my-project
```

Install dependencies

```bash
  npm install
```

Start the server

```bash
  npm run start
```

  
## Usage/Examples

```javascript
import Component from 'my-project'

function App() {
  return <Component />
}
```

  
## Running Tests

To run tests, run the following command

```bash
  npm run test
```

  
## API Reference

#### Get all items

```http
  GET /api/items
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

#### Get item

```http
  GET /api/items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.

  
## Tech Stack

**Client:** React, Redux, TailwindCSS

**Server:** Node, Express

  
## Feedback

If you have any feedback, please reach out to us at fake@fake.com

  
## Authors

- [@katherinepeterson](https://www.github.com/katherinepeterson)

  
## Appendix

Any additional information goes here

  
## FAQ

#### Question 1

Answer 1

#### Question 2

Answer 2

  