# tf-caf-terratest-common

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

Terratest support utitities and test runners supporting CAF framework terraform modules auto tests in pipeline
1. To keep infra tests code DRY and composable, reusable functions extracted into this dedicated repo and to be included by TF modules tests
2. Tests are configuration driven, to be reusable and aggregatable for higher level customer project specific integration testing
3. Configuration is shared across infra deployment (terraform) and infra test (terratest) automation
3. Automated pipelines friendly. Configuration switches by OS env vars


By default test suite pointed to <Module Repo>/examples and expects configuration variables in test.tfvars
```
make check
```

To point to customer project specific terraform code
```
DSO_INFRA_TEST_CONFIG_FOLDER=/projects/abc/ make check
```

To override default test.tfvars
```
DSO_INFRA_TEST_CONFIG_FOLDER=/projects/abc/ DSO_INFRA_TEST_CONFIG_TFVAR_FILENAME=project.tfvars make check
```

tests and inidvidual test stages can be skipped
Example:
<tf module repo>
examples/
   ecs_example/
       main.tf
	   test.tfvars
   eks_example/
		main.tf
		test.tfvars

To skip a test when run test for a module
```
DSO_INFRA_TEST_SKIP_TEST_<name of test TF folder> make check
```
```
DSO_INFRA_TEST_SKIP_TEST_ecs_example make check
```

To disable selected stage(s) of the test
```
SKIP_teardown_test_eks_example=y make check
```

TODO - to pickup from multi tfvars to align to any project file naming convention

## Reusing test impl for high env post deployment regression testing

We want reuse same test implementation for a TF module development and for regression testing of a project that includes the module, probably among multi other ones
Not every test can be reused. As a low level primitive TF module, like "DNS record" has to have an extensive test fixture.
We solve it by introducing naming convention for GoLang tests. Those safe to be composed/reused from higher level pipelines have TestComposable prefix in GoLang test name.
Example:

```
=== ECS-Application-module/tests/testimpl.go ===
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	...
	assert.Equal(t, ctx.TestConfig.(*ThisTFModuleConfig).dockerImage, getAWSEcsAPI().FargateApp(appArn).Container().ImageName)

}
====
```


## Diagrams

[Overview](doc/Overview.svg)


### local development

To test amendments to the terratest helper before those committed to github, use GoLang "replace". Example
```
module github.com/nexient-llc/tf-aws-module-private_dns_namespace

go 1.20

replace github.com/nexient-llc/tf-caf-terratest-common => /Home/user/CAF/NOT_CHECKED_IN_YET/tf-caf-terratest-common

require (
	github.com/nexient-llc/tf-caf-terratest-common v0.0.0-00010101000000-000000000000
)
```

### GoLang

To use "github.com/nexient-llc" private repository when develop or run GoLang code

```
go env -w GOPRIVATE='github.com/nexient-llc/'
```

### Pipeline integration

For unattended CiCd pipelines to preauthenhhicate Github.

for https auth:
```
git config --add --global url."https://oauth2:$GITHUB_PTA_TOKEN@github.com/".insteadOf "https://github.com/"
```
for ssh auth:
```
git config --add --global url."ssh://git@github.com/".insteadOf "https://github.com/"
```

### Prerequisites

- [asdf](https://github.com/asdf-vm/asdf) used for tool version management
- [make](https://www.gnu.org/software/make/) used for automating various functions of the repo
- [repo](https://android.googlesource.com/tools/repo) used to pull in all components to create the full repo template

### Repo Init

Run the following commands to prep repo and enable all `Makefile` commands to run

```shell
asdf plugin add terraform
asdf plugin add tflint
asdf plugin add golang
asdf plugin add golangci-lint
asdf plugin add nodejs
asdf plugin add opa
asdf plugin add conftest
asdf plugin add pre-commit
asdf plugin add terragrunt

asdf install
```

## Pre-Commit hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to terraform, golang and common linting tasks. There are no custom hooks added.

`commitlint` hook enforces commit message in certain format. The commit contains the following structural elements, to communicate intent to the consumers of your commit messages:

- **fix**: a commit of the type `fix` patches a bug in your codebase (this correlates with PATCH in Semantic Versioning).
- **feat**: a commit of the type `feat` introduces a new feature to the codebase (this correlates with MINOR in Semantic Versioning).
- **BREAKING CHANGE**: a commit that has a footer `BREAKING CHANGE:`, or appends a `!` after the type/scope, introduces a breaking API change (correlating with MAJOR in Semantic Versioning). A BREAKING CHANGE can be part of commits of any type.
footers other than BREAKING CHANGE: <description> may be provided and follow a convention similar to git trailer format.
- **build**: a commit of the type `build` adds changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- **chore**: a commit of the type `chore` adds changes that don't modify src or test files
- **ci**: a commit of the type `ci` adds changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- **docs**: a commit of the type `docs` adds documentation only changes
- **perf**: a commit of the type `perf` adds code change that improves performance
- **refactor**: a commit of the type `refactor` adds code change that neither fixes a bug nor adds a feature
- **revert**: a commit of the type `revert` reverts a previous commit
- **style**: a commit of the type `style` adds code changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **test**: a commit of the type `test` adds missing tests or correcting existing tests

Base configuration used for this project is [commitlint-config-conventional (based on the Angular convention)](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional#type-enum)

If you are a developer using vscode, [this](https://marketplace.visualstudio.com/items?itemName=joshbolduc.commitlint) plugin may be helpful.

`detect-secrets-hook` prevents new secrets from being introduced into the baseline. TODO: INSERT DOC LINK ABOUT HOOKS

In order for `pre-commit` hooks to work properly

- You need to have the pre-commit package manager installed. [Here](https://pre-commit.com/#install) are the installation instructions.
- `pre-commit` would install all the hooks when commit message is added by default except for `commitlint` hook. `commitlint` hook would need to be installed manually using the command below

```
pre-commit install --hook-type commit-msg
```

## To run local quality check

1. For development/enhancements to this module locally, you'll need to install all of its components. This is controlled by the `configure` target in the project's [`Makefile`](./Makefile). Before you can run `configure`, familiarize yourself with the variables in the `Makefile` and ensure they're pointing to the right places.

```
make configure
```

This adds in several files and directories that are ignored by `git`. They expose many new Make targets.

2. The first target you care about is `check`.
If `make check` target is successful, developer is good to commit the code to git repo.

`make check` target

- runs `terraform commands` to `lint`,`validate` and `plan` terraform code.
- runs `conftests`. `conftests` make sure `policy` checks are successful.
- runs `terratest`. This is integration test suit.
- runs `opa` tests
