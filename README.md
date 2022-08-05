# Pre-Release

The intention of this is app is to eventually allow the synchronization of git repositories between GitHub, GitLab, BitBucket, and others.

# How To Use

* Download appropriate binary for your OS
* Generate GitLab token with X permissions
* Generate GitHub token with X permissions
* `git2git --github-token 1234 --gitlab-token 5678


# Supported Flags

* --github-token
* --gitlab-token

# Latest Release
Initial working release of git2git. Copies repositories from GitLab to GitHub.

Some caveats:

Only copies repositories that aren't empty.
Only copies repositories that don't have same-named repositories in GitHub.
Only copies git history. Issues and other metadata aren't copied.

# Milestones

- [X] Copy all owned, private repositories from GitLab to GitHub
- [ ] Copy all owned, private repositories from GitHub to GitLab
- [ ] Synchronize all owned, private repositories between GitHub and GitLab
- [ ] Synchronize latest changes from repositories to the next
- [ ] Interactive mode where repositories can be selected
- [ ] Add support for additional flags (archived, public, private, ... repositories)
- [ ] Add support for other git hosting environments
- [ ] Schedule Syncronization
- [ ] Copy issues w/ content from one project to another
- [ ] Zip up repositories and copy to backup location (Google Drive, etc)
