#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

// In a real distribution, we would detect OS and arch to pick the right binary
// For this MVP, we assume the binary is in the package root or standard location
// Or we invoke the go binary directly if it's a dev link.

// Determine binary path (simplified for demo)
const binName = os.platform() === 'win32' ? 'envault.exe' : 'envault';
const binPath = path.join(__dirname, '..', 'bin', binName);

const child = spawn(binPath, process.argv.slice(2), {
    stdio: 'inherit'
});

child.on('exit', (code) => {
    process.exit(code);
});
