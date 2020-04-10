# Logging through pipes

This functionality is needed when a parent process starts a child process and both their logs have to be collected in the parent process.

## The parent process:

```
part, _ := pipes.NewParentPart(marshalizer)
profileReader, logsWriter := part.GetChildPipes()

command = exec.Command("child.bin")
command.ExtraFiles = []*os.File{
		...,
		profileReader,
		logsWriter,
}

part.StartLoop()
```

`StartLoop` will continuously read log lines from the child  (pipe `logsWriter`) on a separate goroutine. Furthermore, the parent part also forwards log profile changes to the child process (through pipe `profileReader`).

Note that the parent is responsible to call `logger.NotifyProfileChange()` when it applies a new log profile (whether by sole choice or when instructed by a logviewer).

## The child process

```
profileReader := os.NewFile(42, "/proc/self/fd/42")
logsWriter := os.NewFile(43, "/proc/self/fd/43")
part := pipes.NewChildPart(profileReader, logsWriter, marshalizer)
err := part.StartLoop()
```

The child has to aquire the provided pipes, create its part of the logging dialogue and then call `StartLoop`.
The child part is automatically registered as observer to the global default `LogOutputSubject`, which means that it gets notified on each log write from any of the loggers in the process. When notified, the child part simply forwards the message (the serialized log line) to its parent, through pipe `logsWriter`. 

Furthermore, the child part listens for eventual log profile changes on the pipe `profileReader`. Any profile change is applied immediately.

### Capturing child text output (stdout, stderr)

On the parent's side:

```
command := exec.Command("child.bin")
childStdout, _:= command.StdoutPipe()
arwenStderr, _ := command.StderrPipe()
parentPart.ContinuouslyReadTextualOutput(childStdout, arwenStderr, "child-tag")
```

`Stdout` will be logged with `trace` level, while `stderr` with `error` level.
