const jwt = require('jsonwebtoken');

// Kong JWT credentials
const key = 'user-key';
const secret = 'user-secret';

function generateToken(options = {}) {
    const payload = {
        iss: key,          // issuer ต้องตรงกับ key ที่สร้างใน Kong
        exp: Math.floor(Date.now() / 1000) + (options.expireIn || 60 * 60 * 24), // default 24 hours
        iat: Math.floor(Date.now() / 1000),
        kid: key,          // key id ต้องตรงกับ key ที่สร้างใน Kong
        ...options.claims  // additional claims
    };

    return jwt.sign(payload, secret);
}

// ถ้าเรียกโดยตรง
if (require.main === module) {
    const token = generateToken();
    console.log('\nGenerated JWT Token:');
    console.log(token);
    console.log('\nCurl command for testing:');
    console.log(`curl -i -X GET http://localhost:8000/api/v1/products -H "Authorization: Bearer ${token}"`);
}

module.exports = generateToken;
