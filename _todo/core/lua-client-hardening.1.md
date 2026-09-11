# Lua client hardening

- Area: transpiler
- Phase: multi-language Dreego
- Goal: finish the supported browser Lua contract with isolation, safe emission, and source-accurate diagnostics.
- Acceptance: adversarial scripts cannot escape generated component scope, emitted JavaScript remains syntactically safe, failures identify the original `.dreego` source range, and black-box browser tests cover every supported Lua construct.
- Depends on: shipped Lua browser foundation
