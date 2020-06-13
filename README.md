# bugalert

_Work in progress_

A git like interface for working with bugs and issues.

```
Work with bugs and issues

Usage:
  bugalert [flags]
  bugalert [command]

Available Commands:
  help        Help about any command
  history     show bug history
  list        list all bugs/issues
  show        show bug/issue details
  version     Print version

Flags:
  -h, --help   help for bugalert

Use "bugalert [command] --help" for more information about a command.

```

Reads config from `~/.bugalert.yml`:

```yaml
url: "https://bugzilla.example.com"
api_key: XXXXXXYYYYYYYXXXXXXYYYYYYXXXXYYXYXYXY
default_user: you@example.com
default_filter: "status:OPEN|ASSIGNED"
```

The API key for accessing bugzilla can be obtained from `API Keys` tab in
`preferences` from your bugzilla site.

### Basic Usage

List all the bugs, no filtering, and for all users.

```
bugalert list --nofilter --all
```

Override the default filter in the config file, and also the default user.

```
bugalert list --filter "status:CLOSED|NEEDINFO" --user somebody@example.com
```

Further more information about a particular bug can be obtained with the `show`
and `history` commands.

```
bugalert show 100
bugalert show --comments 100
bugalert history 100
```

100 in the example is the bug ID.


### TODO
- [ ] Support to update bugs/add comments
- [ ] Support to edit a bug and push
- [ ] Add support for Github issues
- [ ] Caching of issues locally
