#!/usr/bin/env node

const { install } = require('./postinstall');

// This script is a wrapper around the postinstall script to allow for standalone execution
if (require.main === module) {
  install();
}
