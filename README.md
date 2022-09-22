# Pre-Release

The intention of this is app is to eventually allow the synchronization of git repositories between GitHub, GitLab, BitBucket, and others.

# How To Use

1. Download appropriate binary for your OS
2. Generate GitLab token with permissions:

![image](https://user-images.githubusercontent.com/34039134/191862977-0cedcde6-d730-4e70-a3e4-59ddc2ebe6a9.png)

3. Generate GitHub token with permissions:

![image](https://user-images.githubusercontent.com/34039134/191863169-865c6a2b-5a05-4d26-a0f1-64a1bcd172a3.png)

4. `git2git --target-token 1234 --origin-token 5678 --target gitlab --origin github`


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
- [X] Synchronize latest changes from repositories to the next
- [X] Interactive mode where repositories can be selected
- [ ] Add support for additional flags (archived, public, private, ... repositories)
- [ ] Add support for other git hosting environments
- [ ] Schedule Synchronization
- [ ] Copy issues w/ content from one project to another
- [ ] Zip up repositories and copy to backup location (Google Drive, etc)
