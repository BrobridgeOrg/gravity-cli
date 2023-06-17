#!/bin/sh

# Prepare product
./gravity-cli product create accounts --desc="testing product" --enabled --schema=./schema_test.json

# Add rule to product
./gravity-cli product ruleset add accounts accountCreated --enabled --event=accountCreated --method=create --handler=./handler_test.js --schema=./schema_test.json --pk=id
