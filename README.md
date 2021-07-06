# 🦇 &nbsp;Gothic

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/jrapoport/gothic/test?style=flat-square) [![Go Report Card](https://goreportcard.com/badge/github.com/jrapoport/gothic?style=flat-square&)](https://goreportcard.com/report/github.com/jrapoport/gothic) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jrapoport/gothic?style=flat-square) [![GitHub](https://img.shields.io/github/license/jrapoport/gothic?style=flat-square)](https://github.com/jrapoport/gothic/blob/master/LICENSE)

[![Buy Me A Coffee](https://img.shields.io/badge/buy%20me%20a%20coffee-☕-6F4E37?style=flat-square)](https://www.buymeacoffee.com/jrapoport)

Gothic is a user registration and authentication microservice written in Go. It's based on OAuth2 and JWT and will
handle user signup, authentication and custom user data.

## A Complete Rewrite

~85-90% complete.

This project was originally forked from
[Netlify's GoTrue](https://github.com/netlify/gotrue). It has been completely rewritten from the ground 
up and literally has no code in common with the original fork. The basic idea was to support most of the existing 
functionality of GoTrue while expanding support for additional external providers and functionality.

GroTrue relied on its own custom implementations for oauth, mail templates, smtp, jwt, etc. These have all been
replaced with mature external libraries. Details below.

### gRPC / gRPC Web support

**This is currently in progress**. The external grpc api should not be considered stable and will 
likely as things come online.

### Configuration
Gothic uses [viper](https://github.com/spf13/viper) for configuration support.

Gothic supports config files in env, yaml and json formats in addition to env vars.

Please see the [example config](https://github.com/jrapoport/gothic/blob/master/example.env) or 
[test configs](https://github.com/jrapoport/gothic/blob/master/config/testdata) for examples and 
a summary of the currently supported configuration settings.

### Database
Gothic uses [gorm](https://gorm.io/) for database support.

### OAuth

Gothic uses [goth](https://github.com/markbates/goth) for external oauth providers now. Now we support everything
that [goth](https://github.com/markbates/goth) supports.

### Email

Gothic uses [hermes](https://github.com/matcornic/hermes/) for email templates now. 

### SMTP

Gothic uses [go-simple-mail](https://github.com/xhit/go-simple-mail/) for smtp server support. 

### JWT 

Gothic uses [jwt v4](https://github.com/dgrijalva/jwt-go) for smtp server support.


## Project History

The original purpose was to adopt newer, more developer friendly technologies like
[Gorm](https://gorm.io/), [gRPC](https://grpc.io/), and [gRPC Web](https://github.com/grpc/grpc-web); newer versions of
critical libraries like [JWT v4](https://github.com/dgrijalva/jwt-go); and migrate away from older libraries that are
deprecated with [security flaws](https://github.com/gobuffalo/uuid).

These changes allow for advances like self-contained database migration, expanded database driver support (e.g.,
PostgreSQL), and gRPC support. Broadly speaking, they are intended to make it easier to modify and use the microservice
outside of Netlify tool chain, and in a more active development environment.

While the Netlify team did a good job with GoTrue, their use in production means they cannot easily adopt these kinds of
significant changes. In many cases, they will likely never make them given the impacts to their tooling, deployment, and
production systems — which makes perfect sense for their situation.

I'd like to thank Netlify team for their hard work on the original version of this microservice.
