#!/usr/bin/env node

const { execSync } = require('child_process');

// Get arguments from the command line
const args = process.argv.slice(2).join(' ');

// Build the command to execute
const command = `go test ${args} ./...`;

console.log(`Running tests: ${command}`);

try {
  // Execute the command
  execSync(command, { stdio: 'inherit' });
  console.log('✅ Tests passed!');
} catch (error) {
  console.error(`❌ Tests failed: ${error.message}`);
  process.exit(1);
}
