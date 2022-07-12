# quick-task-creator
Simple and minimalistic CLI tool to create trello tasks from my shell

## Geting started

Compile and rename the executable to ```trello```.

Authorize with trello
```
$ trello authorize
```

Configure
```
$ trello configure
```
For the tool to function properly, it needs to know the trello list id.
You can determine the list id by appending ```.json``` to the board address in your browser
and searching the resulting output.

Watch out for a section that looks like this
```
list": {
    "id": "5c1e833954388516a4ef2980",
    "name": "Backlog"
}
```

The tool will store its configuration and authentication token inside ```${HOME}/.quick-task-creator```
