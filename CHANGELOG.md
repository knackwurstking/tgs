# CHANGELOG

## v2.0.0 [2025-07-28]

**Breaking Changes**:

- Refactored all bot commands to use extensions with golang build tags
- Restructured configuration system - each extension now has separate config files
- Renamed "tgs-server" to "tgs"
- Removed the `-c` option from the tgs command
- Changed extension interface - removed Register and Targets methods

**Added**:

- New `pgvis` extension for pg-vis integration
    - `/pgvissignup` command for user registration
    - API key generation and management
    - User creation and profile linking
    - Deep linking support with UUID validation
- Callback query handling system
- Reply callbacks system for interactive commands
- User and chat target validation
- Extension start method for initialization
- Makefile for build automation
- Support for inline keyboards and reply markups
- User name generation with first/last name fallbacks
- Logging for all received updates
- Exit code handling for external commands

**Changed**:

- Updated UI library from v2.0.0 to v4.3.0
- Replaced slog with standard log package
- User names are no longer required to be unique
- Always trim API keys before setting
- Set report caller to true for logging
- Improved error handling with client notifications
- Enhanced debug logging throughout the system
- Updated example configurations for all extensions
- Moved all commands to extensions directory structure
- Renamed configuration files from `.config` to `.yaml`

**Fixed**:

- Nil pointer issues in various components
- Wrong API key assignment logic
- User creation process
- Configuration paths for extensions
- Reply message ID handling in command checks
- Go build tags syntax
- Missing function body errors
- Callback query target validation
- Extension configuration loading
- Exit code handling for pg-vis commands
- Message formatting and escaping

**Removed**:

- Name generator dependency (no longer needed)
- Old monolithic command structure
- Unused debug logs and comments
- Legacy configuration options

## v1.2.0 [2025-04-10]

**Added**:

- IPv6 address support in `/ip` bot command

**Fixed**:

- Invalid character handling in IPv6 URLs

## v1.1.0 [2025-02-11]

**Updated**:

- Updated UI library to [v2.0.0](https://github.com/knackwurstking/ui/tree/dev?tab=readme-ov-file)

**Added**:

- Logging for new users

## v1.0.0 [2024-12-11]

**Initial Release**:

- Commands:
    - `/ip` - Get current IP address information
    - `/stats` - Display system statistics
    - `/journallist` - List available journal entries
    - `/journal` - Access journal functionality
    - `/opmangalist` - List available manga
    - `/opmanga` - Access manga functionality
