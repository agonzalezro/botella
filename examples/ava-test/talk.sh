#!/bin/sh

read line

echo "Hi!"
echo ""
echo "You can use environment variable for sharing secrets with me. For example, in your ava.yml you specified a var called KEY with value: $KEY"
echo ""
echo "I also received this as input: $line"
echo ""
echo "It is just json, parse it, for example, with jq I can know that the body was: $(echo $line|jq .body)"
