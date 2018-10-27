# git-templates

## Intro

These are just my personal git hooks -- I provide them here in a feeble attempt to
ever so slightly lower the overall entropy of the universe (IOW: I hope they make
someone's life a little easier).

## Install

To install, do something like:

    mkdir -p ~/.git-templates/hooks
    cp $PreCommitHookScript ~/.git-templates/hooks/pre-commit
    git config --global init.templatedir '~/.git-templates'
    
Now, all new repos will get a copy of the hook scripts from your template directory.

To install into existing repos, simply copy the appropriate hook script into the repo's
`.git/hooks` directory.

For all the gory details, take a look at https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks

## Description

Currently there's just one hook script here:

* pre-commit -- Never (accidentally) commit an ELF binary ever again!!

The included pre-commit hook runs `file` against each file staged for commit and,
if any file looks like an ELF binary executable, you'll be prompted for what to
do next. The possible options are:

* Attempt to automatically fix the issue
* Abort the commit altogether -- allowing you to fix things manually
* Continue on and actually commit the binary

If you choose to auto-fix the issue, the following steps will be performed:

1. Each detected binary file will be unstaged from the current commit
2. If the unstaged file's basename is not mentioned in the repository's `.gitignore` file, it will be appended there.
3. If any changes are made to the repo's `.gitignore`, that file will be staged for the current commit.
4. If all of the above are successful, the commit will proceed -- otherwise and error will be emitted and the commit will be aborted.

If you really do plan to commit executables to the repo you can skip these checks (and
associated prompt) by setting the `git config` variable `precommit.allow-binaries` to
a `true` value, like so:

    git config --bool precommit.allow-binaries 1

Of course, you could always just remove the hook script from your repo too.
