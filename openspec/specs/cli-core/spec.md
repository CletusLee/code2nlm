## ADDED Requirements

### Requirement: CLI Execution & Flags
The system SHALL expose all arguments natively using Cobra bindings to correctly parse `--input`, `--output`, `--max-sources`, `--max-words`, `--ignore-file`, and `--strategy`.

#### Scenario: Successful Parameter Parsing
- **WHEN** user executes tool with custom parameters
- **THEN** the application binds them cleanly and launches execution lifecycle

### Requirement: Strict Capacity Validations
The system SHALL immediately abort returning a non-zero exit code if required files exceed `--max-sources` constraints.

#### Scenario: Boundary Limits Exceeded
- **WHEN** the overall scanned word count specifies 85 output files required and user enforces a limit of 50
- **THEN** tool terminates immediately, emits a FATAL error specifying the delta cleanly to stderr, and exits with a non-zero status
