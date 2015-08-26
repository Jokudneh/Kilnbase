// +build tools

// the build flag is needed to avoid side-effects of init() functions in tool dependencies.

// Go modules support developer tools (commands) as dependencies.
// For example, your project might require a tool to help with code generation, or to lint/vet your code for correctness.
