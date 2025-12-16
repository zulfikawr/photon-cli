#!/usr/bin/env node

/**
 * Build script for bitrim
 * Builds Go binaries for all supported platforms
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const platforms = [
  { os: 'windows', arch: 'amd64', name: 'bitrim-windows-amd64.exe' },
  { os: 'darwin', arch: 'amd64', name: 'bitrim-darwin-amd64' },
  { os: 'linux', arch: 'amd64', name: 'bitrim-linux-amd64' },
];

const projectRoot = path.join(__dirname, '..');
const goRoot = path.join(projectRoot, '..');
const binDir = path.join(projectRoot, 'bin');

console.log('\n🔨 Building bitrim binaries for all platforms...\n');

// Ensure bin directory exists
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

let successful = 0;
let failed = 0;

platforms.forEach((platform) => {
  const outputPath = path.join(binDir, platform.name);
  
  try {
    console.log(`📦 Building ${platform.name}...`);
    
    // Set environment variables for cross-compilation
    const env = Object.assign({}, process.env, {
      GOOS: platform.os,
      GOARCH: platform.arch,
    });
    
    // Build the binary
    execSync(`go build -o "${outputPath}" main.go`, {
      cwd: goRoot,
      env: env,
      stdio: 'inherit',
    });
    
    console.log(`✅ Built ${platform.name}\n`);
    successful++;
  } catch (error) {
    console.log(`❌ Failed to build ${platform.name}`);
    console.log(`   Error: ${error.message}\n`);
    failed++;
  }
});

// Summary
console.log(`${'='.repeat(50)}`);
console.log(`Build Summary:`);
console.log(`  ✅ Successful: ${successful}/${platforms.length}`);
console.log(`  ❌ Failed: ${failed}/${platforms.length}`);

if (failed === 0) {
  console.log(`\n✅ All binaries built successfully!`);
  console.log(`📁 Binaries location: ${binDir}\n`);
  process.exit(0);
} else {
  console.log(`\n❌ Some builds failed. Check errors above.\n`);
  process.exit(1);
}
