#!/usr/bin/env node

const { execSync } = require('child_process');

// Get arguments from the command line
const args = process.argv.slice(2).join(' ');

// Build the command to execute
const command = `go build -ldflags="-s -w" -o ./bin/bitrim ${args} ./main.go`;

console.log(`Building bitrim: ${command}`);

try {
  // Execute the command
  execSync(command, { stdio: 'inherit' });
  console.log('✅ Build successful!');
} catch (error) {
  console.error(`❌ Build failed: ${error.message}`);
  process.exit(1);
}
