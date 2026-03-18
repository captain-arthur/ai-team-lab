import http from 'k6/http';
import { check, sleep } from 'k6';

const TARGET_URL = __ENV.TARGET_URL || 'https://test.k6.io/';
const MODE = (__ENV.MODE || 'rps').toLowerCase(); // rps | vu
const THINK_TIME_S = parseFloat(__ENV.THINK_TIME_S || '0');

function thresholds() {
  const p95 = parseInt(__ENV.SLO_P95_MS || '800', 10);
  const failRate = parseFloat(__ENV.SLO_FAIL_RATE || '0.01');
  return {
    http_req_failed: [`rate<${failRate}`],
    http_req_duration: [`p(95)<${p95}`],
    checks: ['rate>0.99'],
  };
}

function scenarios() {
  if (MODE === 'vu') {
    return {
      ramp_vus: {
        executor: 'ramping-vus',
        startVUs: 1,
        stages: [
          { duration: '10s', target: 1 },
          { duration: '20s', target: 10 },
          { duration: '10s', target: 2 },
        ],
        gracefulRampDown: '5s',
      },
    };
  }

  return {
    ramp_rps: {
      executor: 'ramping-arrival-rate',
      timeUnit: '1s',
      preAllocatedVUs: 20,
      maxVUs: 200,
      stages: [
        { duration: '10s', target: 5 },
        { duration: '20s', target: 50 },
        { duration: '10s', target: 10 },
      ],
    },
  };
}

export const options = {
  scenarios: scenarios(),
  thresholds: thresholds(),
};

export default function () {
  const res = http.get(TARGET_URL);
  check(res, { 'status is 200': (r) => r.status === 200 });
  if (THINK_TIME_S > 0) sleep(THINK_TIME_S);
}
