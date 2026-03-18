import http from 'k6/http';
import { check } from 'k6';

// CAT Adapter 표준 입력(환경변수)
const TARGET_URL = __ENV.TARGET_URL;
const MODE = (__ENV.MODE || 'arrival').toLowerCase(); // arrival | vus

const DURATION = __ENV.DURATION || '30s';
const TARGET_RPS = parseInt(__ENV.TARGET_RPS || '50', 10);
const VUS = parseInt(__ENV.VUS || '10', 10);

const SLO_P95_MS = parseInt(__ENV.SLO_P95_MS || '400', 10);
const SLO_FAIL_RATE = parseFloat(__ENV.SLO_FAIL_RATE || '0.01');

if (!TARGET_URL) {
  throw new Error('TARGET_URL is required');
}

function scenarios() {
  if (MODE === 'vus') {
    return {
      vus_steady: {
        executor: 'constant-vus',
        vus: VUS,
        duration: DURATION,
      },
    };
  }

  return {
    arrival_steady: {
      executor: 'constant-arrival-rate',
      rate: TARGET_RPS,
      timeUnit: '1s',
      duration: DURATION,
      preAllocatedVUs: 20,
      maxVUs: 200,
    },
  };
}

export const options = {
  scenarios: scenarios(),
  thresholds: {
    http_req_duration: [`p(95)<${SLO_P95_MS}`],
    http_req_failed: [`rate<${SLO_FAIL_RATE}`],
    checks: ['rate>0.99'],
  },
};

export default function () {
  const res = http.get(TARGET_URL);
  check(res, { 'status 200': (r) => r.status === 200 });
}

