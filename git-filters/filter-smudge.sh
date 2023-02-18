#!/bin/sh
sed -i ''  's/"password_hash": *"[^"]*"/"password_hash": "{password_hash}"/g' ../definitions.json