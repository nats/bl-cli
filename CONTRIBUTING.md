# Contributing to bl-cli

<!-- Non emacs users, feel free to update the toc by hand. -->
<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Contributing to bl-cli](#contributing-to-bl-cli)
    - [Issues](#issues)
        - [Reporting an Issue](#reporting-an-issue)
        - [Issue Lifecycle](#issue-lifecycle)
    - [Developing](#developing)
        - [Go environment](#go-environment)
        - [Docker](#docker)
        - [Testing](#testing)
            - [Writing Tests](#writing-tests)
              - [Unit tests](#unit-tests)
              - [Integration tests](#integration-tests)
            - [`go-binarylane` mocks](#go-binarylane-mocks)
            - [Build Scripts](#build-scripts)

<!-- markdown-toc end -->

**First:** if you're unsure or afraid of _anything_, just ask
or submit the issue or pull request anyways. You won't be yelled at for
giving your best effort. The worst that can happen is that you'll be
politely asked to change something. We appreciate all contributions!

For those folks who want a bit more guidance on the best way to
contribute to the project, read on. Addressing the points below
lets us merge or address your contributions quickly.

## Issues

### Reporting an Issue

* Make sure you test against the latest released version. It is possible
  we already fixed the bug you're experiencing.

* If you experienced a panic, please create a [gist](https://gist.github.com)
  of the *entire* generated crash log for us to look at. Double check
  no sensitive items were in the log.

* Respond as promptly as possible to any questions made by the _bl-cli_
  team to your issue. Stale issues will be closed.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a _bl-cli_ collaborator.
   Categorization is done via labels. For example, bugs are marked as "bugs".

3. Unless it is critical, the issue is left for a period of time (sometimes
   many weeks), giving outside contributors a chance to address the issue.

4. The issue is addressed in a pull request. The issue will be
   referenced in commit message(s) so that the code that fixes it is clearly
   linked.

5. The issue is closed. Sometimes, valid issues will be closed to keep
   the issue tracker clean. The issue is still indexed and available for
   future viewers, or can be re-opened if necessary.

## Pull Requests

Pull requests must always be opened from a fork of `bl-cli`, even if you have
commit rights to the repository so that all contributors follow the same process.

## Developing

`bl-cli` has `make` commands for most tooling in the `Makefile`. Run `make`
or `make help` for a list of available commands with descriptions.

## Documenting

`bl-cli` commands have two kinds of documentation: the short synopsis, that shows in the command lists, and the long description, that shows in the `--help` message for a specific command. In `commands/*.go` you'll see these two things being defined frequently, often as different arguments in `CmdBuilderWithDocs`. Here are some guidelines to keep in mind when writing these helpful texts:

- Go uses "quotes" for single-line strings and \``backticks`\` for multi-line strings.
- Programmatic elements, such as command and flag names, should be surrounded by backticks.
- To feature a backtick inside a multiline string, use this sequence of characters for each backtick:

  ```
  ` + "`" + `
  ```
- It's good practice to create string variables to store text that gets repeated.
- Flags and short command synopses do not need complete sentences in their descriptions and should not end in punctuation
- Command abstracts, on the other hand, are considered full-text documentation and should use proper English
- Write short command descriptions from the perspective of the user trying to do something (e.g. "List all database clusters") vs. what the command does (e.g. "This command retrieves a list of all database clusters").
- Avoid the passive voice ("When a tag is provided, access is granted") and use the active voice ("Entering a tag provides access")
- Be helpful when users have to enter a input that is from a list of possible values. Give examples, list the possible values inline (if the list is relatively short), or point them to a command that can list the possible values for them.


### Go environment

The minimal version of Golang for `bl-cli` is 1.14. `bl-cli` uses [Go
modules](https://github.com/golang/go/wiki/Modules) for dependency
management [with vendoring](https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away).
Please run `make vendor` after any dependency modifications.

Be sure to run `go fmt` on your code before submitting a pull request.

### Docker

You can build bl-cli in a local Docker container via `make docker_build`.

### Testing

Run the tests locally via `make test`, or on Travis CI by pushing a branch to your fork
on github.

#### Writing Tests

In bl-cli, we have two kinds of tests: unit tests and integration tests. Both are run with Go's
built-in `go test` tool.

##### Unit tests

Unit tests live in the `_test.go` files. The bulk of these tests live in the `commands` package,
and exist to ensure that a CLI command produces an expected output. For each unit test, we
typically rely on an accompanying mocked API call. These mocks are generated via `gomock`, and
can be set to return different values from the API calls to test how our commands behave when
given different inputs.

Writing a unit test for a new command typically looks like this,

1. Write your new command.
2. If your new command depends on a mocked `go-binarylane` call, generate a mock for it. See
[the section below](#go-binarylane-mocks) about regenerating mocks to learn how to do this.
3. Use your new mocks to stub out the API call, and write a test case. We use
`github.com/stretchr/testify/assert` for our assertions. Test cases typically look like the following:
    ```go
    func TestMyNewCommand(t *testing.T) {
        // Use the `withTestClient` helper to access our tets config, as well as the go-binarylane API
        // mocks.
        withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
            // Mock the go-binarylane API call.
            tm.myNewCommandMock.EXPECT().Get("some-value").Return("some-other-value")

            // Optionally add a CLI argument.
            config.Args = append(config.Args, "some-value")

            // Optionally add a CLI flag.
            config.Doit.Set(config.NS, "--my-flag", "some-value")

            // Execute your command.
            err := RunMyNewCommand(config)

            // Add assertions to check if your test passes.
            assert.NoError(t, err)
        })
    }
    ```

##### Integration tests

Integration tests live under the top-level `integration` directory. These tests exist to ensure
that an invocation of a command though this CLI produces the expected output. These tests use a
mocked HTTP client, but run the actual compiled bl binary.

Writing an integration test typically looks like this,

1. Write your new command.
2. Mock the API's JSON response that your command depends on.
3. Execute your command using `exec.Command` on the test CLI binary.
4. Add assertions to check the output from the CLI command.

Rather than providing an example here, please have a look at the [`integration/account_test.go`](/integration/account_test.go)
file to see what an integration test typically looks like.

#### `go-binarylane` mocks

To upgrade `go-binarylane`, run `make upgrade_binarylane`. This will:

* Get the latest release of `go-binarylane`, and update the go.mod and go.sum files accordingly.
* Tidy and vendor the modules that `bl-cli` depends on.
* Run `mockgen` to regenerate the mocks we use in the unit test suite.

#### Build Scripts

If you modify the build scripts, you can use `make shellcheck` to
check your changes. You'll need to install [shellcheck](https://github.com/koalaman/shellcheck)
to do so. Travis also runs shellcheck.
