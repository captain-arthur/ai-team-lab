import http from 'k6/http';
import { check, group, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'https://test.k6.io';
const THINK_TIME_S = parseFloat(__ENV.THINK_TIME_S || '0.5');

export const options = {
  vus: 2,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<800'],
    checks: ['rate>0.99'],
  },
};

export default function () {
  group('step1_home', () => {
    const res = http.get(`${BASE_URL}/`);
    check(res, { 'home 200': (r) => r.status === 200 });
  });
  sleep(THINK_TIME_S);

  group('step2_pi', () => {
    const res = http.get(`${BASE_URL}/pi.php?decimals=3`);
    check(res, { 'pi 200': (r) => r.status === 200 });
  });
  sleep(THINK_TIME_S);

  group('step3_contacts', () => {
    const res = http.get(`${BASE_URL}/contacts.php`);
    check(res, { 'contacts 200': (r) => r.status === 200 });
  });
}
