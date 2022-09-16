# Pre-Release

The intention of this is app is to eventually allow the synchronization of git repositories between GitHub, GitLab, BitBucket, and others.

# How To Use

* Download appropriate binary for your OS
* Generate GitLab token with X permissions
* Generate GitHub token with X permissions
* `git2git --target-token 1234 --origin-token 5678 --target gitlab --origin github


# Supported Flags

* --target-token
* --origin-token
* --target (gitlab,github)
* --origin (gitlab,github)

# Latest Release
Supports sync of all branches for (github -> gitlab) (gitlab -> github) (github -> github) (gitlab -> gitlab).

Some caveats:

Only copies repositories that aren't empty.
Only copies repositories that don't have same-named repositories in .
Only copies git history. Issues and other metadata aren't copied.

## Warning: Do not try this if you have repositories in the target that have the same name as the ones you're trying to sync.

# Milestones

- [X] Copy all owned, private repositories from GitLab to GitHub
- [X] Copy all owned, private repositories from GitHub to GitLab
- [X] Copy all owned, private repositories from GitLab to GitLab
- [X] Copy all owned, private repositories from GitHub to GitHub
- [ ] Synchronize all owned, private repositories between GitHub and GitLab
- [ ] Synchronize latest changes from repositories to the next
- [ ] Interactive mode where repositories can be selected
- [ ] Add support for additional flags (archived, public, private, ... repositories)
- [ ] Add support for other git hosting environments
- [ ] Schedule Syncronization
- [ ] Copy issues w/ content from one project to another
- [ ] Zip up repositories and copy to backup location (Google Drive, etc)
