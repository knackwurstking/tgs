# CHANGELOG

## v2.1.1 [2025-08-02]

**Changed**:

- Updated all references from "pg-vis" to "pg-press"
- Changed "Vis" to "Presse" (German)

**Fixed**:

- Corrected German terminology and naming conventions

## v2.1.0 [2025-08-02]

**Breaking Changes**:

- Replaced `pgvis` extension with `pgpress` extension
- Changed command from `/pgvisregister` to `/pgpressregister`
- Updated deep linking URLs to use `pgpressregister` instead of `pgvisregister`

**Added**:

- New `pgpress` extension for PG-Vis Server integration (replaces `pgvis`)
    - `/pgpressregister` command for user registration with API key generation
    - Deep linking support with UUID validation for secure registration
    - User creation and profile linking with external pg-vis command integration
    - Inline keyboard for direct link to PG-Vis Server login page
    - User existence checking and API key management
    - Support for both new user creation and existing user API key retrieval

**Changed**:

- Updated telegram-bot-api dependency to v5.6.0 (custom fork)
- Added `pgpress` to default extension build tags in Makefile (replacing `pgvis`)
- Enhanced error handling for external command execution with proper exit code handling
- Improved user name generation with fallback to first/last name combination

**Fixed**:

- Telegram bot API library compatibility issues with custom fork integration
- User API key generation and validation workflow
- External command execution error handling with specific exit codes
- Deep link validation and expiration handling

**Removed**:

- `pgvis` extension and its `/pgvisregister` command (replaced by `pgpress`)
