const http = require('http');
const generateToken = require('./generate-token');

async function testToken(token, endpoint = '/api/v1/products') {
    return new Promise((resolve, reject) => {
        const options = {
            hostname: 'localhost',
            port: 8000,
            path: endpoint,
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        };

        const req = http.request(options, (res) => {
            let data = '';
            
            console.log('\nResponse Status:', res.statusCode);
            console.log('Response Headers:', res.headers);

            res.on('data', (chunk) => {
                data += chunk;
            });

            res.on('end', () => {
                try {
                    if (data) {
                        const jsonData = JSON.parse(data);
                        console.log('\nResponse Body:', JSON.stringify(jsonData, null, 2));
                    }
                    resolve({
                        status: res.statusCode,
                        headers: res.headers,
                        body: data ? JSON.parse(data) : null
                    });
                } catch (error) {
                    console.log('\nRaw Response:', data);
                    resolve({
                        status: res.statusCode,
                        headers: res.headers,
                        body: data
                    });
                }
            });
        });

        req.on('error', (error) => {
            console.error('\nError:', error.message);
            reject(error);
        });

        req.end();
    });
}

async function runTests() {
    console.log('Running JWT token tests...\n');

    // Test 1: Valid token
    console.log('Test 1: Testing with valid token');
    const validToken = generateToken();
    await testToken(validToken);

    // Test 2: Expired token
    console.log('\nTest 2: Testing with expired token');
    const expiredToken = generateToken({ expireIn: -3600 }); // expired 1 hour ago
    await testToken(expiredToken);

    // Test 3: Invalid signature
    console.log('\nTest 3: Testing with invalid signature');
    const invalidToken = validToken.slice(0, -1) + 'X';
    await testToken(invalidToken);

    // Test 4: Different endpoints
    console.log('\nTest 4: Testing different endpoints');
    const endpoints = [
        '/api/v1/products',
        '/api/v1/orders',
        '/api/v1/auth/profile'
    ];
    
    for (const endpoint of endpoints) {
        console.log(`\nTesting endpoint: ${endpoint}`);
        await testToken(validToken, endpoint);
    }
}

// Run tests
runTests().catch(console.error);
