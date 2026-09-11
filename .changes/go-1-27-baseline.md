---
version: minor
---

- Breaking: require Go 1.27 or newer for Dreego development, generated applications, and builds.
- Changed: adopt idiomatic language, standard-library, runtime, and tooling improvements from Go 1.23 through Go 1.27.
- Security: confine plugin client modules and documentation reads to their declared filesystem roots and cap repeated HTTP header values.
