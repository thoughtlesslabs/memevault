# Changelog

All notable changes to this project will be documented in this file.

## [v1.2.1] - 2026-01-14

### Security
- **Secret Injection Fix**: Fixed a vulnerability where `memevault get` output could be spoofed by secrets containing newlines.
- **Input Validation**: `memevault set` now rejects keys with non-alphanumeric characters and values with newlines/control characters.
- **Runtime Hardening**: `memevault run` now actively filters out invalid keys from the loaded vault before injecting them into the environment, preventing potential attacks from malicious vault files.

### Changed
- `memevault get` now outputs values in quoted format (e.g., `KEY="Value"`) for better safety and shell compatibility.
