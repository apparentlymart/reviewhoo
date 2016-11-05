# reviewhoo

`reviewhoo` (pronounced "review-hoo") is a simple tool that uses git history
to suggest a few possible code reviewers for a particular branch.

It's far from an exact science, because such a tool can never be perfect.
Instead, the goal is to make a "best effort" with the following goals in mind:

* Share the review workload somewhat-evenly between those who are well-suited
  to review a particular change.

* Each change ideally wants two different kinds of reviewer: one who is
  an expert on the part of the code being changed, and one that is expert
  on the technology or infrastructure involved but is less familiar with the
  particular code being changed.

`reviewhoo` works by first identifying the files and directories that are
most affected by a given change. It then analyzes which other authors have
made a lot of commits in those areas recently, and chooses one of them
at random.

It then extracts the file extensions of all of the affected files and looks
for the most common ones as an approximation of what programming languages
are used primarily by the change. We then finally seek out the most prolific
recent committers on files with these extensions and again choose one at
random for each extension.

# Usage

Here's what a typical run might look like for a commit that changes both
JavaScript and SCSS files:

```
$ reviewhoo

reviewhoo suggests the following reviewers:

* **Sofía Ramirez** has changed these files recently
* **Abelina Brown** has made lots of .js changes recently
* **Garrett McCarthy** has made lots of .scss changes recently

```

This command expects to be run from a git work tree where a topic branch is
currently selected. By default it will analyze all of the commits that are
on the current branch but not on `master`.

The output is intended to be Markdown-ish so that it can be pasted verbatim
into a Github Pull Request comment, or similar. Of course, in that case you
may wish to manually replace the author names with github usernames so you
can 'at-mention' the people who you'd like to do a review.

*Sofía* was selected by looking at a tally of recent commits by author in the
files most modified (in terms of lines changed) by the current branch.

*Abelina* was selected by looking at the recent commit history for commits
that changed files whose names end with `.js`, because the current branch
contains lots of changes to `.js` files.

*Garrett* was selected similarly by looking for commits that change files
ending with `.scss`, for the same reason.

The branch might also have made small changes to ancillary files like
`package.json`, but these are not considered because they are minor changes
(in terms of line count) compared to the `.js` and `.scss` files.

`reviewhoo` only consults the most recent two months of history when looking
for candidate reviewers. The goal of this limit is to only request reviews
from those who are currently active contributors to a codebase, under the
assumption that people who have drifted away from a codebase will be less
effective as reviewers. Of course, as a user you are free to ignore the
suggestions and directly ask a particular person to review if you like.

The flip-side of this two month rule is that when developers leave a team
they can still be identified as reviewers for up to two months after their
last commit. In future we might extend the tool to be able to take as input
a file describing the currently-active reviewers (e.g. an export of the
current members of a Github team) but this is not yet implemented; if you're
not happy with `reviewhoo`'s initial suggestions, you can of course just run
it again to get another set, assuming that there are other reviewers that
qualify.

## More Advanced Usage

Sometimes you have a topic branch that is intended to be merged into *another*
topic branch. In this case you can override the default of comparing the
current branch to `master`, by specifying another reference branch on the
command line:

```
$ reviewhoo base-branch
```

You can also, if you need to, ask for a set of reviewers for a branch other
than the one currently on, by specifying additionally the branch name where
the new changes can be found:

```
$ reviewhoo base-branch change-branch
```

The defaults for the base branch and the change branch, used when not specified
on the command line, are `master` and `HEAD` respectively.

## JSON Output

If desired, you can use `reviewhoo` to gather data for presentation using
another tool, by requesting JSON output:

```
$ reviewhoo -json
```

In the JSON output the author names and reasons are presented as separate
strings so that, for example, you can look up the author names in a mapping
table to find corresponding Github usernames.
