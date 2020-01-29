# Milestone Closer GitHub Action

This action closes Issue Milestones in a GitHub repository that are past-due and
have zero open issues.

## Example usage

```yaml
on:
  schedule:
    # Cron syntax has five fields separated by a space, and each field represents a unit of time.
    #
    #        ┌───────────── minute (0 - 59)
    #        │ ┌───────────── hour (0 - 23)
    #        │ │ ┌───────────── day of the month (1 - 31)
    #        │ │ │ ┌───────────── month (1 - 12 or JAN-DEC)
    #        │ │ │ │ ┌───────────── day of the week (0 - 6 or SUN-SAT)
    #        │ │ │ │ │
    #        │ │ │ │ │
    #        │ │ │ │ │
    #        * * * * *
    # * is a special character in YAML so you have to quote this string
    - cron:  '0 0 * * SUN' # midnight every Sunday

jobs:
  milestones:
    runs-on: ubuntu-latest
    name: Milestone Closer
    steps:
    - name: Close Milestones
      uses: MrGossett/github-action-milestone-closer@v1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
