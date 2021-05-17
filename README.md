# git-new
Creates a git branch based on a Jira issue ID.
It takes a Jira issue ID and converts it to a common name and checkouts a new branch with that name.

Example:

```bash
git new FM-2806
```

Will change to a branch named:

```bash
2806-translate-account-cancellation-suspension-journey
```

## Installation
### Getting the JIRA token
Create a new token in https://id.atlassian.com/manage-profile/security/api-tokens

### Configure the ENV vars

```bash
export JIRA_SERVER=https://your.jira.server
export PATH=$PATH:$HOME/bin
export JIRA_USER=your@user.com
export JIRA_TOKEN=yoursupersecrettoken
```

### Copy the file and give permissions
*NOTE*: `ln -s` might not work in some shells, better if you copy the file.

```bash
mkdir -p ~/bin/
cp git-new ~/bin/
chmod +x ~/bin/git-new
```
