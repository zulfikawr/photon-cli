#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const RELEASE_TAG = 'v1.0.0';
const GITHUB_OWNER = 'zulfikawr';
const GITHUB_REPO = 'photon-cli';
const GITHUB_API = `https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/tags/${RELEASE_TAG}`;

// Determine platform and architecture
function getPlatformInfo() {
  const platform = process.platform;
  const arch = process.arch;

  if (platform === 'win32') {
    return {
      platform: 'windows',
      arch: arch === 'x64' ? 'amd64' : arch,
      filename: 'photon-cli-windows-amd64.exe',
      executable: false // Windows doesn't need chmod
    };
  } else if (platform === 'darwin') {
    return {
      platform: 'darwin',
      arch: arch === 'x64' ? 'amd64' : arch,
      filename: 'photon-cli-darwin-amd64',
      executable: true
    };
  } else if (platform === 'linux') {
    return {
      platform: 'linux',
      arch: arch === 'x64' ? 'amd64' : arch,
      filename: 'photon-cli-linux-amd64',
      executable: true
    };
  }

  throw new Error(`Unsupported platform: ${platform}`);
}

// Download binary from GitHub release
function downloadBinary(downloadUrl, outputPath) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(outputPath);
    
    https.get(downloadUrl, (response) => {
      // Handle redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        downloadBinary(response.headers.location, outputPath)
          .then(resolve)
          .catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download binary: HTTP ${response.statusCode}`));
        return;
      }

      response.pipe(file);

      file.on('finish', () => {
        file.close();
        resolve();
      });

      file.on('error', (err) => {
        fs.unlink(outputPath, () => {}); // Delete the file if error occurs
        reject(err);
      });
    }).on('error', (err) => {
      fs.unlink(outputPath, () => {}); // Delete the file if error occurs
      reject(err);
    });
  });
}

// Get download URL from GitHub release
function getDownloadUrl(filename) {
  return new Promise((resolve, reject) => {
    https.get(GITHUB_API, { headers: { 'User-Agent': 'photon-cli-npm' } }, (response) => {
      let data = '';

      response.on('data', (chunk) => {
        data += chunk;
      });

      response.on('end', () => {
        try {
          const release = JSON.parse(data);
          const asset = release.assets.find((a) => a.name === filename);

          if (!asset) {
            reject(new Error(`Binary not found in release: ${filename}`));
            return;
          }

          resolve(asset.browser_download_url);
        } catch (err) {
          reject(new Error(`Failed to parse GitHub API response: ${err.message}`));
        }
      });
    }).on('error', reject);
  });
}

async function install() {
  try {
    const platformInfo = getPlatformInfo();
    const binDir = path.join(__dirname, '..', 'bin');
    const binPath = path.join(binDir, process.platform === 'win32' ? 'photon-cli.exe' : 'photon-cli');

    // Create bin directory if it doesn't exist
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }

    console.log(`📥 Downloading photon-cli for ${platformInfo.platform}-${platformInfo.arch}...`);

    const downloadUrl = await getDownloadUrl(platformInfo.filename);
    await downloadBinary(downloadUrl, binPath);

    // Make binary executable on Unix-like systems
    if (platformInfo.executable) {
      fs.chmodSync(binPath, 0o755);
    }

    console.log(`✅ Successfully installed photon-cli to ${binPath}`);
    console.log(`🚀 Try it out: photon-cli --help`);
  } catch (error) {
    console.error(`❌ Installation failed: ${error.message}`);
    process.exit(1);
  }
}

// Only run if this is being executed directly (during npm install)
if (require.main === module) {
  install();
}

module.exports = { install, getPlatformInfo, downloadBinary, getDownloadUrl };
