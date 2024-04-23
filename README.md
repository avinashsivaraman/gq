# gq
Simple CLI to ask questions


## Building the project

```
> go build
```

## Running the project

```
> ./gq

> ./gq "Seattle" -q "What is the weather like today?"

> echo "Seattle" | ./gq -q "What is the weather like today?"
```

## Providing An API Key

```
> touch ~/.config/gq/.gq.yaml
```

Provide the below api_key in the file
```
api_key: a28f29de-be82-4d03-afd3-a94f36414a7d
```
