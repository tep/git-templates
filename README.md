# git-templates

## Intro

These are just my personal git hooks -- I provide them here in a feeble attempt
to ever so slightly lower the overall entropy of the universe (IOW: I hope they
make someone's life a little easier).

## Install

To install, do something like:

    mkdir -p ~/.git-templates/hooks
    cp $PreCommitHookScript ~/.git-templates/hooks/pre-commit
    git config --global init.templatedir '~/.git-templates'
    
Now, all new repos will get a copy of the hook scripts from your template
directory.

To install into existing repos, simply copy the appropriate hook script into
the repo's `.git/hooks` directory.

For all the gory details, take a look at
https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks

## Description

Currently there's just one hook script here:

### pre-commit -- Never (accidentally) commit an ELF binary ever again!!

On Linux, there's no simple way to globally declare that you want `git` to
ignore every binary file; unlike Windows, there's no conventional filename
extension to indicate that a file is executable. To remedy this, I've written
a quick-n-dirty pre-commit hook script to detect when/if an executable has
accidentally been staged for commit.

NOTE: This hook script won't work on Windows -- but, then again, you don't
      really need it on Windows, do you.

The included pre-commit hook runs the Linux `file` command against each file
staged for commit and, if any of them look like an ELF binary executable, you'll
be prompted for what to do next. The possible options are:

* Attempt to automatically fix the issue
* Abort the commit altogether -- allowing you to fix things manually
* Continue on and actually commit the binary

If you choose to auto-fix the issue, the following steps will be performed:

1. Each detected binary file will be unstaged from the current commit
2. If the unstaged file's basename is not mentioned explicitly in the
   repository's `.gitignore` file, it will be appended there.
3. If any changes are made to the repo's `.gitignore`, that file will be staged
   for the current commit.
4. If all of the above are successful, the commit will proceed -- otherwise an
   error will be emitted and the commit will be aborted.

If you choose to manually abort the commit, you'll probably want to do all
or most of the above steps yourself.

If you choose to commit the binary -- and plan to continue doing so in the
future -- you can skip these checks (and associated prompt) by setting the
`git config` variable `precommit.allow-binaries` to a `true` value, like so:

    git config --bool precommit.allow-binaries 1

Of course, you could always just remove the hook script from your repo too.
