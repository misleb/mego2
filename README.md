# Mego2

A web application built with Go and WebAssembly.

## Demo

The application is currently deployed and accessible at:

**[https://mego2-efe371cd64de.herokuapp.com/](https://mego2-efe371cd64de.herokuapp.com/)**

## Project Structure

- `app/` - Client-side application (compiled to WebAssembly)
- `server/` - Backend server and API endpoints
- `shared/` - Shared types and utilities
- `web/` - Static web assets

## Roadmap

- [ ] Allow Google users to set a local password
- [ ] Create a DB column "type" and enforce based on mappings for a model when using orm.Update (prevent invalid column names at compile time)
- [ ] Allow anyone to change password or name
- [ ] Enforce password strength frontend and backend
- [ ] Expand test coverage for shared packages
- [ ] Implement admin dashboard for managing content

## Development

See the `Makefile` for available build and development commands.

