# tmux-too-young: The Very Special tmux Session Opener...

## About

`tmux-too-young` is tmux session opener which includes features for git, git worktrees and [tmuxp](https://github.com/tmux-python/tmuxp).

It makes it effortless to launch and switch between tmux sessions when your projects contain a git repo.

If your git repo uses git worktrees then each individual worktree is offered as a different session.

And, if your repo contains a `.tmuxp.yaml` file, tmuxp will be used to launch your session, meaning your tmuxp session file will be used to configure the windows and panels.

## Usage

### How do I run it?

To open `tmux-too-young`, run:

```bash
./tmux-too-young
```

See 'Are there any tips in terms of integrating it into my workflow?' below for tips on making it less of a verbose pain to run.

### What gets listed?

Having run `tmux-too-young`, you are presented with a list of potential sessions. Some may already exist, some may not.

Each session represents either:

| Item | Displayed | Example |
|-------------|-----------|---------|
| A git repo  | The name of the root directory | `tmux-too-young` |
| A git worktree within a repo (if the repo supports worktrees) | The name of the root directory followed by name of branch | `tmux-too-young -> main` |

You can enter a search term to filter the list of potential sessions, and navigate to select one.

### What happens when I select an item?

* If a session with the name does not exist: a new session is created and attached to.
* If a session with the name exists: the session is attached to.

### The intro mentions `tmuxp` support.....?

Yes, as a bonus if a `.tmuxp.yaml` file exists in the root, `tmuxp` will be used to launch the session, meaning your tmuxp session file will be used to determine its windows, panels etc.

Including a `.tmuxp.yaml` file in your repo and lauching it via `tmux-too-young` makes for a fantastic dev experience. Recommended!

### Can I pass in a search term via the command line?

Yes. You can pass an argument to `tmux-too-young`. This string will be used as your initial search term.

If the search term returns only one result, that session will be automattically created/attached to.

This is useful for a couple reasons:

1. Sometimes its just nice to write the search term whilst you are already in the act of writing the command.
1. Doing so means it becomes part of your command history, meaning you can easily recall it using all the usual shell tricks.

### Are there any tips in terms of integrating it into my workflow?

There are.

1. It's recommended you create an alias to call it. Whilst `tmux-too-young` makes for a wonderful pun, its a bit much to type.
1. You can add it as a key binding within your tmux config, allowing you to trigger it with a couple key presses from anywhere within tmux.

## Development

### Running

```
go run tmux-too-young
```

### "Deploying"

Note, this will change once you formalise things.

For the time being, this is how you get it into the bin directory so you can use it personally across your machines.

```
go build tmux-too-young
mv tmux-too-young ~/dotfiles/bin/bin/
```
