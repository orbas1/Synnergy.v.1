const crypto = require('crypto');

// Derive a deterministic session token from a private key. In a
// production implementation this would interface with the user's wallet
// or secure key store. The private key is expected to be supplied via an
// environment variable and never hard-coded.
function authenticate(privateKey) {
  if (!privateKey) {
    throw new Error('DESKTOP_SHELL_PRIVATE_KEY not set');
  }

  return crypto
    .createHmac('sha256', privateKey)
    .update('synnergy-session')
    .digest('hex');
}

module.exports = { authenticate };
