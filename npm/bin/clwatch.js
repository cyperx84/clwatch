#!/usr/bin/env node
"use strict";

const { execFileSync } = require("child_process");
const path = require("path");

const bin = path.join(__dirname, process.platform === "win32" ? "clwatch.exe" : "clwatch");

try {
  execFileSync(bin, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  process.exit(err.status ?? 1);
}
