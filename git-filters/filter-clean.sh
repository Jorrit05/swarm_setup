#!/bin/sh
sed -i ''  's/"password_hash": *"[^"]*"/"password_hash": "{password_hash}"/g' ../definitions3.json