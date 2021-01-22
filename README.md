# Food Delivery

Source code of Food Delivery app.

## Get Started

Requirements:

- [Docker](https://www.docker.com/get-started)

List all available commands:

```
make all
```

### Development

Run in Development mode:

```
make dev
```

`Cmd+C` won't stop dev mode, use belowed command to stop dev mode:

```
make dev_stop
```

If there was any change in project's configuration, then run this command to rebuild dev environment:

```
make dev_build
```

Visit `localhost:8081` to see db's GUI (phpMyAdmin).

### Deployment

To build executable file in local:

```
make build
```

To build docker image for deployment:

```
make build_image
```

The above way will build image with default name: `food-delivery`. However, if you want to change the name then run this:

```
make build_image NAME={{custom-image-name}}
```

## Project Structure

This project follow [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Contribution Guidelines

### Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0) to make commit messages.
