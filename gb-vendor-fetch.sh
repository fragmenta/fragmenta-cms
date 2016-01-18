#!/usr/bin/env bash

gb vendor fetch github.com/fragmenta/assets

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/lib/authorise

# fetching recursive dependency github.com/fragmenta/auth
# fetching recursive dependency github.com/fragmenta/fragmenta-cms/src/users
# fetching recursive dependency github.com/fragmenta/fragmenta-cms/src/lib/status
# fetching recursive dependency github.com/fragmenta/model
# fetching recursive dependency github.com/fragmenta/query
# fetching recursive dependency github.com/fragmenta/router
# fetching recursive dependency github.com/fragmenta/server
# fetching recursive dependency github.com/fragmenta/view/helpers
# fetching recursive dependency github.com/go-sql-driver/mysql
# fetching recursive dependency github.com/kennygrant/sanitize
# fetching recursive dependency github.com/lib/pq
# fetching recursive dependency golang.org/x/crypto/bcrypt
# fetching recursive dependency golang.org/x/crypto/blowfish
# fetching recursive dependency golang.org/x/net/html

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/pages/actions

# fetching recursive dependency github.com/fragmenta/fragmenta-cms/src/pages
# fetching recursive dependency github.com/fragmenta/fragmenta-cms/src/posts
# fetching recursive dependency github.com/fragmenta/view

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/posts/actions

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/tags/actions

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/images
# fetching recursive dependency github.com/bamiaux/rez

gb vendor fetch github.com/fragmenta/fragmenta-cms/src/app

gb vendor fetch github.com/sendgrid/sendgrid-go
# fetching recursive dependency github.com/sendgrid/smtpapi-go

# after some time has passed, run:
# gb vendor update --all

