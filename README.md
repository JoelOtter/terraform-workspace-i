# `terraform-workspace-i`

This is an interactive tool for switching Terraform workspaces. It's based on
[git branch-i](https://github.com/JoelOtter/git-branch-i), which does the same
thing for Git branches.

## Installation

Install it using Go, like so:

```sh
go install github.com/JoelOtter/terraform-workspace-i@latest
```

Ensure your Go directory is on your system path. You might want to alias this
tool to something easier to type - I use `tfw`.

## Usage

* Workspaces can be navigated using the arrow keys, j and k, Pg Up/Down, or
Ctrl+N and Ctrl+P.
* Select a workspace with the return key.
* Delete a workspace with the delete or backspace key, and use y/n to confirm.
* Exit with Escape or Ctrl+C.
