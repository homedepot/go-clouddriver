# Contributions Welcome

First off, thank you for considering contributing to Go Clouddriver!

If you're just looking for quick feedback for an idea or proposal, feel free to open an
[issue](https://github.com/homedepot/go-clouddriver/issues/new).

Follow the [contribution workflow](#contribution-workflow) for submitting your
changes to the Go Clouddriver codebase.

## Contribution Workflow

The Go Clouddriver repository uses the “fork-and-pull” development model. Follow these steps if
you want to merge your changes to the Go Clouddriver repository:

1. Within your fork of 
   [Go Clouddriver](https://github.com/homedepot/go-clouddriver), create a
   branch for your contribution. Use a meaningful name.
2. Create your contribution, meeting all
   [contribution quality standards](#contribution-quality-standards)
3. [Create a pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)
   against the master branch of the Go Clouddriver repository in the Homedepot org. Make sure to set the base repository as
   `homedepot/go-clouddriver`.
4. Add a reviewer to your pull request. Work with your reviewer to address any comments and obtain an approval.
   To update your pull request amend existing commits whenever applicable and
   then push the new changes to your pull request branch.
5. Once the pull request is approved, one of the [maintainers](MAINTAINERS.md) will merge it.

## Contribution Quality Standards

Your contribution needs to meet the following standards:

- Separate each **logical change** into its own commit.
- Add a descriptive message for each commit. Follow
  [commit message best practices](https://github.com/erlang/otp/wiki/writing-good-commit-messages).
- Document your pull requests. Include the reasoning behind each change, and
  the testing done.
- Acknowledge the [Apache 2.0 license](LICENSE).
