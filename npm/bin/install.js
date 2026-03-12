#!/usr/bin/env node
"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");
const zlib = require("zlib");

const VERSION = require("../package.json").version;
const REPO = "cyperx84/clwatch";

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  const osMap = { darwin: "darwin", linux: "linux", win32: "windows" };
  const archMap = { arm64: "arm64", x64: "amd64" };

  const goos = osMap[platform];
  const goarch = archMap[arch];

  if (!goos || !goarch) {
    throw new Error(`Unsupported platform: ${platform}/${arch}`);
  }

  return { goos, goarch };
}

function download(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        return download(res.headers.location).then(resolve, reject);
      }
      if (res.statusCode !== 200) {
        return reject(new Error(`HTTP ${res.statusCode} for ${url}`));
      }
      const chunks = [];
      res.on("data", (chunk) => chunks.push(chunk));
      res.on("end", () => resolve(Buffer.concat(chunks)));
      res.on("error", reject);
    }).on("error", reject);
  });
}

async function main() {
  const { goos, goarch } = getPlatform();
  const archiveName = `clwatch-${VERSION}-${goos}-${goarch}.tar.gz`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${archiveName}`;
  const binDir = path.join(__dirname);
  const binaryName = goos === "windows" ? "clwatch.exe" : "clwatch";
  const binaryPath = path.join(binDir, binaryName);

  console.log(`Downloading clwatch v${VERSION} for ${goos}/${goarch}...`);

  const tarGz = await download(url);

  // Extract: write tar.gz to temp, use tar to extract
  const tmpFile = path.join(os.tmpdir(), archiveName);
  const tmpDir = path.join(os.tmpdir(), `clwatch-extract-${Date.now()}`);
  fs.writeFileSync(tmpFile, tarGz);
  fs.mkdirSync(tmpDir, { recursive: true });

  execSync(`tar -xzf "${tmpFile}" -C "${tmpDir}"`);
  fs.unlinkSync(tmpFile);

  // Find the binary in extracted directory
  const extractedDir = path.join(tmpDir, `clwatch-${VERSION}-${goos}-${goarch}`);
  const extractedBin = path.join(extractedDir, binaryName);

  fs.copyFileSync(extractedBin, binaryPath);
  fs.chmodSync(binaryPath, 0o755);
  fs.rmSync(tmpDir, { recursive: true, force: true });

  console.log(`Installed clwatch to ${binaryPath}`);
}

main().catch((err) => {
  console.error(`Failed to install clwatch: ${err.message}`);
  process.exit(1);
});
