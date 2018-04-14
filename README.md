# JenkinsProxy
A proxy to let webhooks to reach a Jenkins instance running behind a firewall

## PROBLEM

Jenkins is awesome and matchless tool for both CI & CD; but unfortunately its a gold mine if left in wild with wide open access; so, we always want to put it behind a firewall. But when we put it behind firewall then webhooks don't work anymore and no one wants the pull based polling but rather prefer the build to start as soon as there is a commit!

## SOLUTION

This little proxy makes webhooks start working again!

## STATUS

WORK IN PROGRESS
