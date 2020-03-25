# elrond-go-logger
Elrond's logger subsystem written in go

## CLI options

### Logs producer (Elrond Node)

 - `log-level`: comma-separated pairs of (`loggerName`, `logLevel`) 
 - `log-correlation`: option to include correlation elements in the logs
 - `log-logger-name`: option to include logger name in the logs

Example:

```
--log-level="*:INFO,processor:DEBUG" --log-correlation --log-logger-name
```

### Logs viewer

```
--level="*:INFO,processor:DEBUG" --correlation --logger-name
```
