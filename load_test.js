import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

export let options = {
    vus: 250,
    duration: '1m',
    thresholds: {
        http_req_duration: ['p(99)<50'],
        checks: ['rate>0.9999']
    },
};

const BASE_URL = 'http://localhost:8085/api';

let responseTime = new Trend('response_time');
let successRate = new Rate('success_rate');

export default function () {
    let username = `test_user${__VU}`;
    let authPayload = JSON.stringify({ username: username });
    let params = { headers: { 'Content-Type': 'application/json' } };

    let authRes = http.post(`${BASE_URL}/auth`, authPayload, params);
    responseTime.add(authRes.timings.duration);
    
    let authSuccess = check(authRes, {
        'Авторизация успешна': (res) => res.status === 200 && res.json('token') !== null,
    });
    successRate.add(authSuccess);
    
    if (!authSuccess) {
        console.log(`Ошибка авторизации: ${authRes.status} - ${authRes.body}`);
        return;
    }

    let authToken = authRes.json('token');
    let authHeaders = { headers: { 'Authorization': `Bearer ${authToken}` } };

    // Получение информации о пользователе
    let infoRes = http.get(`${BASE_URL}/info`, authHeaders);
    let infoSuccess = check(infoRes, { 'Информация получена': (res) => res.status === 200 });
    successRate.add(infoSuccess);
    if (!infoSuccess) console.log(`Ошибка info: ${infoRes.status} - ${infoRes.body}`);

    let sendPayload = JSON.stringify({ toUser: 'Denis', amount: 1 });
    let sendRes = http.post(`${BASE_URL}/sendCoin`, sendPayload, authHeaders);
    let sendSuccess = check(sendRes, { 'Монеты отправлены': (res) => res.status === 200 });
    successRate.add(sendSuccess);
    if (!sendSuccess) console.log(`Ошибка sendCoin: ${sendRes.status} - ${sendRes.body}`);

    let buyRes = http.get(`${BASE_URL}/buy/pen`, authHeaders);
    let buySuccess = check(buyRes, { 'Покупка прошла успешно': (res) => res.status === 200 });
    successRate.add(buySuccess);
    if (!buySuccess) console.log(`Ошибка buy: ${buyRes.status} - ${buyRes.body}`);


    sleep(1);
}
