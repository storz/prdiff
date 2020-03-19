## prdiff [![GoDoc][godoc-badge]][godoc] [![Go Report Card][go-report-badge]][go-report]

prdiff enumerates GitHub Pull Requests merged into default branch since last release (or between two branches/commits/tags)

[godoc]: https://godoc.org/github.com/storz/prdiff
[godoc-badge]: https://godoc.org/github.com/storz/prdiff?status.svg
[go-report]: https://goreportcard.com/report/github.com/storz/prdiff
[go-report-badge]: https://goreportcard.com/badge/github.com/storz/prdiff

### Installation

```sh
go get -u github.com/storz/prdiff/cmd/prdiff

# if couldn't find any packages because of you are inside $GOPATH
GO111MODULE=on go get -u github.com/storz/prdiff/cmd/prdiff
```

### Usage

```sh
prdiff [[base] head]
prdiff [base...head]
```

#### Arguments

```
base   begin point of diff (default: latest release)
head   end point of diff (default: default branch set on GitHub)
```

#### Options

```
  -o, --owner string    user/organization name (default: repository of current directory)
  -r, --repo string     repository name (default: repository of current directory)
  -t, --target string   specified target branch what you want to get diff (default: default branch set on GitHub)
      --token string    GitHub access token (default: $GITHUB_TOKEN)
```
