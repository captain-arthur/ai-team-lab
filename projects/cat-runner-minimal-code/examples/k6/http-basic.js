import http from 'k6/http';
import { check, sleep } from 'k6';

// CAT runner가 주입하는 값
const TARGET_URL = __ENV.TARGET_URL;
const MODE = (__ENV.MODE || 'arrival').toLowerCase(); // arrival | vus
const DURATION = __ENV.DURATION || '20s';

const TARGET_RPS = parseInt(__ENV.TARGET_RPS || '50', 10);
const VUS = parseInt(__ENV.VUS || '10', 10);

const SLO_P95_MS = parseInt(__ENV.SLO_P95_MS || '400', 10);
const SLO_FAIL_RATE = parseFloat(__ENV.SLO_FAIL_RATE || '0.01');

if (!TARGET_URL) {
  throw new Error('TARGET_URL is required');
}

function scenario() {
  if (MODE === 'vus') {
    return {
      vus_steady: {
        executor: 'constant-vus',
        vus: VUS,
        duration: DURATION,
      },
    };
  }

  // default: arrival-rate
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
  scenarios: scenario(),
  thresholds: {
    http_req_duration: [`p(95)<${SLO_P95_MS}`],
    http_req_failed: [`rate<${SLO_FAIL_RATE}`],
    checks: ['rate>0.99'],
  },
};

export default function () {
  const res = http.get(TARGET_URL);
  check(res, { 'status is 200': (r) => r.status === 200 });
  // think time은 이 최소 구현에서는 사용하지 않는다.
  // sleep(0.01);
  sleep(0);
}

