#!/usr/bin/env node

/**
 * Manual download script for photon-cli binary
 * Run with: node scripts/download.js
 */

const { install } = require('./postinstall.js');

install().catch((err) => {
  console.error('❌ Download failed:', err.message);
  process.exit(1);
});
