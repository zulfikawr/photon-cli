#!/usr/bin/env node

/**
 * Test script for bitrim npm package
 * Validates package structure before publishing
 */

const fs = require('fs');
const path = require('path');

const tests = [];
let passed = 0;
let failed = 0;

// Helper to add test
function test(name, fn) {
  tests.push({ name, fn });
}

// Helper to assert
function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

// Tests
test('package.json exists and is valid', () => {
  const pkgPath = path.join(__dirname, '..', 'package.json');
  assert(fs.existsSync(pkgPath), 'package.json not found');
  
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
  assert(pkg.name === 'bitrim', 'Invalid package name');
  assert(pkg.version, 'Missing version');
  assert(pkg.author, 'Missing author');
  assert(pkg.license, 'Missing license');
});

test('npm package version matches Go CLI version', () => {
  const pkgPath = path.join(__dirname, '..', 'package.json');
  const changelogPath = path.join(__dirname, '..', '..', 'CHANGELOG.md');
  
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
  const changelog = fs.readFileSync(changelogPath, 'utf8');
  
  // Extract version from CHANGELOG.md header (## [X.Y.Z])
  const versionMatch = changelog.match(/## \[([^\]]+)\]/);
  assert(versionMatch, 'Could not find version in CHANGELOG.md');
  
  const goVersion = versionMatch[1];
  assert(pkg.version === goVersion, 
    `npm version (${pkg.version}) does not match Go CLI version (${goVersion}). Keep versions in sync!`
  );
});

test('bin directory exists', () => {
  const binPath = path.join(__dirname, '..', 'bin');
  assert(fs.existsSync(binPath), 'bin directory not found');
  assert(fs.statSync(binPath).isDirectory(), 'bin is not a directory');
});

test('bin/bitrim wrapper exists', () => {
  const binPath = path.join(__dirname, '..', 'bin', 'bitrim');
  assert(fs.existsSync(binPath), 'bin/bitrim wrapper not found');
});

test('bin/bitrim.cmd wrapper exists', () => {
  const cmdPath = path.join(__dirname, '..', 'bin', 'bitrim.cmd');
  assert(fs.existsSync(cmdPath), 'bin/bitrim.cmd wrapper not found');
});

test('scripts directory exists', () => {
  const scriptsPath = path.join(__dirname, '..', 'scripts');
  assert(fs.existsSync(scriptsPath), 'scripts directory not found');
  assert(fs.statSync(scriptsPath).isDirectory(), 'scripts is not a directory');
});

test('postinstall.js exists', () => {
  const postinstallPath = path.join(__dirname, 'postinstall.js');
  assert(fs.existsSync(postinstallPath), 'postinstall.js not found');
});

test('download.js exists', () => {
  const downloadPath = path.join(__dirname, 'download.js');
  assert(fs.existsSync(downloadPath), 'download.js not found');
});

test('README.md exists', () => {
  const readmePath = path.join(__dirname, '..', 'README.md');
  assert(fs.existsSync(readmePath), 'README.md not found');
});

test('postinstall.js is valid JavaScript', () => {
  const postinstallPath = path.join(__dirname, 'postinstall.js');
  const content = fs.readFileSync(postinstallPath, 'utf8');
  assert(content.includes('getPlatformInfo'), 'Missing getPlatformInfo function');
  assert(content.includes('downloadBinary'), 'Missing downloadBinary function');
  assert(content.includes('getDownloadUrl'), 'Missing getDownloadUrl function');
});

test('bin wrappers are not empty', () => {
  const unixPath = path.join(__dirname, '..', 'bin', 'bitrim');
  const cmdPath = path.join(__dirname, '..', 'bin', 'bitrim.cmd');
  
  const unixContent = fs.readFileSync(unixPath, 'utf8');
  const cmdContent = fs.readFileSync(cmdPath, 'utf8');
  
  assert(unixContent.length > 0, 'Unix wrapper is empty');
  assert(cmdContent.length > 0, 'Windows wrapper is empty');
});

// Run all tests
console.log('\n📋 Running bitrim npm package tests...\n');

tests.forEach(({ name, fn }) => {
  try {
    fn();
    console.log(`✅ ${name}`);
    passed++;
  } catch (error) {
    console.log(`❌ ${name}`);
    console.log(`   Error: ${error.message}`);
    failed++;
  }
});

// Summary
console.log(`\n${'-'.repeat(50)}`);
console.log(`Tests passed: ${passed}/${tests.length}`);

if (failed === 0) {
  console.log('✅ All tests passed! Ready to publish.\n');
  process.exit(0);
} else {
  console.log(`❌ ${failed} test(s) failed. Fix issues before publishing.\n`);
  process.exit(1);
}
