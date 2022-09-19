# secgrp

`secgrp` manages the security group on Amazon EC2.

## Install

```
$ git clone https://github.com/rtiwsk/secgrp
$ cd secgrp
$ go build ./cmd/secgrp/
```

## Usage

```
$ secgrp -h
Usage:
  secgrp [options...]

Options:
  -id      Specify the EC2 instance ID.
  -sgid    Specify the security group ID.
  -add     Add a security group to the instance.
  -remove  Remove the security group from the instance.
  -list    List the security group for the instance.

Example:
  $ secgrp -id i-1234567890abcdef -list
  $ secgrp -id i-1234567890abcdef -sgid sg-1234567890abcdefg -add
  $ secgrp -id i-1234567890abcdef -sgid sg-1234567890abcdefg -remove
```
