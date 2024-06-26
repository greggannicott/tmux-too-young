# tmux-too-young: The Very Special tmux Session Opener...

## IMPORTANT

This is still a work in progress with no official first version. Frankly I'm amazed you've found it!

Use it at your own risk.

## About

`tmux-too-young` is tmux session opener which includes features for git, git worktrees and [tmuxp](https://github.com/tmux-python/tmuxp).

It makes effortless work of launching and switching between tmux sessions when your projects contain a git repo or a `.tmuxp.yaml` file.

If your git repo uses git worktrees then each individual worktree is offered as a different session.

And, if your repo contains a `.tmuxp.yaml` file, tmuxp will be used to launch your session, meaning your tmuxp session file will be used to configure the windows and panels.

## Installation

### Homebrew

Install with:

```bash
brew install greggannicott/tools/tmux-too-young
```

Update to the latest version using:

```bash
brew update
brew upgrade tmux-too-young
```

### Other

Looking for an alternative way to install it? Please create an issue and we'll attempt to provide it.

## Setup

The first time you run `tmux-too-young open`, you will be prompted to enter the directories you want it to scan for projects (ie. the directories/sessions that are listed).

Enter a path, press enter and you will be prompted to enter another. Once you have finished entering your paths, provide an empty entry and it will proceed to display a list of projects.

At this point your `~/.tmux-too-young.conf` file will be created.

## Usage

### How do I run it?

To open a tmux session with `tmux-too-young`, run:

```bash
./tmux-too-young open
```

See 'Are there any tips in terms of integrating it into my workflow?' below for tips on making it less of a verbose pain to run.

### What gets listed?

Having run `tmux-too-young open`, you are presented with a list of potential sessions. Some may already exist, some may not.

Each session represents either:

| Item | Displayed | Example |
|-------------|-----------|---------|
| A git repo  | The name of the root directory | `my-project` |
| A git worktree within a repo (if the repo supports worktrees) | The name of the root directory followed by name of branch | `my-project -> my-branch` |
| A directory containing a `.tmuxp.yaml` file | The name of the root directory | `my-project` |
| A directory containing an empty `.tmux-too-young` file | The name of the root directory | `my-project` |

As [fzf](https://github.com/junegunn/fzf) is used to display the results, you can enter a search term to filter the list of potential sessions, and navigate to select one.

### What happens when I select an item?

* If a session with the name does not exist: a new session is created and attached to.
* If a session with the name exists: the session is attached to.

### The intro mentions `tmuxp` support.....?

Yes, not only does it detect directories with `.tmuxp.yaml` files, but as a bonus if a `.tmuxp.yaml` file exists in the root, `tmuxp` will be used to launch the session, meaning your tmuxp session file will be used to determine its windows, panels etc.

Including a `.tmuxp.yaml` file in your repo and lauching it via `tmux-too-young` makes for a fantastic dev experience. Recommended!

### Can I pass in a search term via the command line?

Yes. You can pass a value to `tmux-too-young`. For example:

```bash
./tmux-too-young open --search "telescope main"
```

The string "telescope main" will be used as your initial search term, with results being filtered by that term.

If the search term returns only one result, that session will be automattically created/attached to.

If we had a project called 'telescope' which supported worktrees, `tmux-too-young` would automattically create/attach to the session for the `main` branch.

Being able to enter the search term when calling the app is useful for a couple reasons:

1. Sometimes its just nice to write the search term whilst you are already in the act of writing the command.
1. Doing so means it becomes part of your command history, meaning you can easily recall it using all the usual shell tricks.

### Are there any tips in terms of integrating it into my workflow?

There are.

1. It's recommended you create an alias to call it. Whilst `tmux-too-young` makes for a wonderful pun, its a bit much to type.
1. You can add it as a key binding within your tmux config, allowing you to trigger it with a couple key presses from anywhere within tmux.

## Development

### Running

```
go run tmux-too-young open
```

### Releasing New Version

#### Push Changes

Make sure remote repo contains latest changes:

```bash
git push
```

#### Update Git Tag

Determine the latest version:

```bash
git tag --sort=-version:refname | head -n 1
```

Update the tag with the new version (following sematic version rules):

```bash
git tag v0.0.1
```

#### Trigger New Release

Push your tags to trigger a new release on Github. This will build the new binaries and update theHomebrew repo using a combination of Goreleaser and Github Actions:

```bash
git push --tags
```
