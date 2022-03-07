# logger

this is just a wrapper on [zap](https://github.com/uber-go/zap)

## Integration

either you can import default config and edit it out according to your needs or you can use `Option` struct to create your own config. for example:

- importing default config and then edit out the values according to your needs
```
lOptions := logger.DefaultOptions()
```

- creating your own config
```
lOptions := logger.Options{
	JSONEncoding: true,
}
```

- pass those options inside **Configure** function to initialize the logger.

```
err := logger.Configure(lOptions)
```

and then you can use that logger according to your needs.
## Example

below are the examples which gonna help you to get started with the the integration.

[example](example/)