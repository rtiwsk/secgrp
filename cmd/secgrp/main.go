package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rtiwsk/secgrp/pkg/securitygroup"
)

var (
	flagid     = flag.String("id", "", "")
	flagsgid   = flag.String("sgid", "", "")
	flagadd    = flag.Bool("add", false, "")
	flagremove = flag.Bool("remove", false, "")
	flaglist   = flag.Bool("list", false, "")
)

var help = `Usage: 
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

`

type securityGroup struct {
	GroupID   string `json:"sgid"`
	GroupName string `json:"name"`
}

type operation int

const (
	none operation = iota
	add
	remove
	list
)

func getOperation() operation {
	if *flaglist {
		return list
	}

	if *flagadd {
		return add
	}

	if *flagremove {
		return remove
	}

	return none
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, help)
	}

	flag.Parse()

	op := getOperation()

	if err := run(op, *flagid, *flagsgid); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(op operation, instanceId, securityGroupId string) error {
	if instanceId == "" {
		return errors.New("The EC2 instance ID is required.")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	sg := securitygroup.New(cfg, instanceId)
	ctx := context.Background()

	switch op {
	case add:
		if err := sg.Add(ctx, securityGroupId); err != nil {
			return err
		}
	case remove:
		if err := sg.Remove(ctx, securityGroupId); err != nil {
			return err
		}
	case list:
		groups, err := sg.List(ctx)
		if err != nil {
			return err
		}

		var sgs []securityGroup
		for _, g := range groups {
			sgs = append(sgs, securityGroup{*g.GroupId, *g.GroupName})
		}

		jsonData, _ := json.MarshalIndent(sgs, "", "    ")
		fmt.Println(string(jsonData))
	case none:
		return errors.New("The operation is not specified.")
	}

	return nil
}
