# bugalert

![Go](https://github.com/fossix/bugalert/workflows/Go/badge.svg?branch=master)
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

#### Listing
List all the bugs, no filtering, and for all users.

```
bugalert list --nofilter --all
```

Override the default filter in the config file, and also the default user.

```
bugalert list --filter "status:CLOSED|NEEDINFO" --user somebody@example.com
```

The listing can be ordered by the last changed time using the `--order`
option. The `--limit N` option can limit the listing to `N`. The below command
gets a list of the last 10 bugs that were closed


```
bugalert list --filter "status:CLOSED" --limit 10 --order
```

#### Show & History
More information about a particular bug can be obtained with the `show` and
`history` commands.

```
bugalert show 100
bugalert show --comments 100
bugalert history 100
```

100 in the example is the bug ID.

#### Updating

To add a new comment for bug 12345

```
bugalert comment 12345
```

This will open a editor, where the comment can be typed in and saved. A comment
can also be provided using the `-m` option. To verify that everything is correct
without actually updating the issue, a `--dry-run` flag can be passed.

The previous comment can be quoted when adding a new comment. This comment can
be obtained from the 'show --comments' listing, where the comment ID is shown as
part of the comment header.

```
$ bugalert show --comments 12345
12345  Very critical issue that needs to be fixed
Status: WORKING
Priority: P1
Severity: ship issue
Created on: 04/07/2020
Creator: User1 <user1@example.com>
Assigned to: User2 <user2@example.com>
QA Contact: User1 <user1@example.com>

Problem description here.

[#100] On 04/07/2020, user2@example.com wrote:
    This looks similar to #12344

$ bugalert comment -q 100 -m "Can we mark this as duplicate of #12344?"
```

The command opens an editor for additional inputs, once the editor is closed, a
new comment is added to the bug/issue.

Note: If the message is empty, and nothing is updated in the editor, an empty
comment is posted with only the quote. So make sure you pass a message with `-m`
or enter text in the editor.

### TODO
- [x] Support to update bugs/add comments
- [ ] Caching of issues locally
- [ ] Support to edit a bug and push
- [ ] Add support for Github issues

