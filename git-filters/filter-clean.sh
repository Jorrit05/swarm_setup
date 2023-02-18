#!/bin/bash

# Remove password_hash field from JSON files
sed -i '' '/"password_hash"/d' $1