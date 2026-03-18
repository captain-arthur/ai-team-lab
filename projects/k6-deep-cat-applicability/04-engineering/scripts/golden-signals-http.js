import http from 'k6/http';
import { check, sleep } from 'k6';

// 목표: latency / error rate / throughput + (간접) saturation 신호를 같이 읽는다.
const TARGET_URL = __ENV.TARGET_URL || 'https://test.k6.io/';
const MODE = (__ENV.MODE || 'arrival').toLowerCase(); // arrival | vus

const P95_MS = parseInt(__ENV.SLO_P95_MS || '400', 10);
const FAIL_RATE = parseFloat(__ENV.SLO_FAIL_RATE || '0.01');

function scenarios() {
  if (MODE === 'vus') {
    const VUS = parseInt(__ENV.VUS || '10', 10);
    const DURATION = __ENV.DURATION || '30s';
    const THINK = parseFloat(__ENV.THINK_TIME_S || '0.0');
    return {
      vus_steady: {
        executor: 'constant-vus',
        vus: VUS,
        duration: DURATION,
        exec: 'vuExec',
        tags: { mode: 'vus' },
        env: { THINK_TIME_S: String(THINK) },
      },
    };
  }

  // arrival-rate는 목표 처리량을 "주입"하고, 달성 실패 시 dropped_iterations로 포화 힌트를 준다.
  const RATE = parseInt(__ENV.TARGET_RPS || '50', 10);
  const DURATION = __ENV.DURATION || '30s';
  return {
    arrival_steady: {
      executor: 'constant-arrival-rate',
      rate: RATE,
      timeUnit: '1s',
      duration: DURATION,
      preAllocatedVUs: parseInt(__ENV.PRE_ALLOCATED_VUS || '20', 10),
      maxVUs: parseInt(__ENV.MAX_VUS || '200', 10),
      exec: 'arrivalExec',
      tags: { mode: 'arrival' },
    },
  };
}

export const options = {
  scenarios: scenarios(),
  thresholds: {
    http_req_failed: [`rate<${FAIL_RATE}`],
    http_req_duration: [`p(95)<${P95_MS}`],
    checks: ['rate>0.99'],
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
};

export function vuExec() {
  const think = parseFloat(__ENV.THINK_TIME_S || '0.0');
  const res = http.get(TARGET_URL, { tags: { test: 'golden' } });
  check(res, { 'status 200': (r) => r.status === 200 });
  if (think > 0) sleep(think);
}

export function arrivalExec() {
  const res = http.get(TARGET_URL, { tags: { test: 'golden' } });
  check(res, { 'status 200': (r) => r.status === 200 });
}
