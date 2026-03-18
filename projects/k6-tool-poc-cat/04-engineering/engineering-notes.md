# Engineering Notes: k6 실행 기반 POC (판단용)

**Date:** 2026-03-18  
**k6 버전:** v1.6.1 (`brew install k6`)

## 실행 대상(POC)
- **타깃**: `https://test.k6.io/`  
- **의미**: 클러스터/ingress 적용 전에, k6의 **실행 모델·지표·판정 메커니즘**을 분리해서 확인하기 위한 최소 타깃

## 공통으로 본 SLI / SLO(POC 최소)
- **SLI**
  - latency: `http_req_duration`의 `p(95)`
  - error rate: `http_req_failed`의 `rate`
  - throughput(참고): `http_reqs`의 `rate`
- **SLO(=threshold)**: 위 SLI에 대한 조건을 주고, k6가 **PASS/FAIL을 종료 코드로 내는지**를 확인

---

## Test 1 — 단일 HTTP 요청 반복 (기본 계측/집계가 “성립”하는지)
- **목적**: k6가 요청→metric→집계(p95/failed rate)→threshold 평가까지 **기본 루프를 정상 수행**하는지 확인
- **검증하려는 것**
  - 기본 HTTP 요청이 SLI로 집계되는가(p95/failed/rate)
  - threshold가 “판정”으로 출력되는가
- **왜 필요한가**: 이게 깨지면 이후의 시나리오/부하 모델 논의는 의미가 없음(도구의 최소 작동 검증)

- **스크립트**: `scripts/test1-single-http.js`
- **실행**

```bash
k6 run --summary-export results/test1-summary.json scripts/test1-single-http.js | tee results/test1-output.txt
```

- **관찰(발췌)**: PASS  
  - `p(95)=189.83ms`, `failed=0.00%`, `http_reqs rate=4.684798/s`
- **해석**
  - k6는 “요청 결과”를 즉시 `http_req_duration`, `http_req_failed`로 기록하고, 종료 시점에 p95/비율을 집계해 threshold를 평가한다.
- **의미(CAT 관점)**
  - k6는 **클라이언트 관측 SLI**를 만들고, 이를 **SLO 게이트**로 바꾸는 기본 기능이 성립한다.

---

## Test 2 — 3-step 사용자 흐름 (다단계 시나리오 표현 + think time의 의미)
- **목적**: “사용자 여정” 형태(2~3 step)를 k6가 **스크립트로 자연스럽게 표현**하고, 그 결과가 SLI로 해석 가능한지 확인
- **검증하려는 것**
  - 다단계 요청이 한 테스트로 묶여 실행/집계되는가
  - `sleep`(think time)이 iteration/처리량에 어떤 영향을 주는지
- **왜 필요한가**: CAT에서 ingress 테스트는 “한 번의 요청”이 아니라 **연속 동작**(페이지/API 체인)인 경우가 많음

- **스크립트**: `scripts/test2-user-flow.js` (home → pi → contacts)
- **실행(think time 포함)**

```bash
THINK_TIME_S=0.5 k6 run --summary-export results/test2-summary.json scripts/test2-user-flow.js | tee results/test2-output.txt
```

- **관찰(발췌)**: PASS  
  - `p(95)=184.59ms`, `failed=0.00%`, `http_reqs rate=7.4493/s`
  - `iteration_duration avg=1.6s` (step 사이 `sleep(0.5)` 포함)
- **해석**
  - `http_req_duration`는 “요청 단위 지연”이고, `iteration_duration`은 “여정(여러 요청+sleep 포함) 한 번”의 시간이다.
  - think time은 요청 자체를 느리게 만들지 않지만, VU가 같은 시간에 수행하는 iteration 수를 줄여 **달성 처리량**을 떨어뜨린다.
- **의미(CAT 관점)**
  - k6는 단일 SLI만 보는 도구가 아니라, “요청 지연”과 “여정 시간”을 분리해 읽을 수 있다.
  - CAT에서 “현실적 사용자”를 흉내 내면 throughput이 줄 수 있으므로, **SLO를 latency만으로 걸지/처리량도 포함할지**가 명확해야 한다.

---

## Test 3 — 부하 모델/판정/think time 실험 (CAT에서 가장 중요한 경계 확인)
- **목적**: CAT 시나리오 정의에서 결정적인 3가지를 분리해서 확인
  - (1) 부하 모델 축: **VU 기반 vs arrival-rate 기반**
  - (2) SLO 게이트: threshold가 **실제로 FAIL을 만들고 종료 코드가 바뀌는지**
  - (3) think time이 “처리량”을 어떻게 바꾸는지
- **검증하려는 것**
  - 같은 타깃에서 executor 선택이 결과 해석(특히 throughput)에 어떤 의미를 가지는지
  - threshold가 “보고서”가 아니라 “판정 장치”로 동작하는지
- **왜 필요한가**: CAT에서 k6를 쓸 때 가장 흔한 실패는 “부하 정의가 모호해 SLO 해석이 흔들리는 것”

- **스크립트**: `scripts/test3-load-pattern.js`

### 3-A) VU 기반 램프(ramping-vus)
```bash
MODE=vu THINK_TIME_S=0 k6 run --summary-export results/test3-vu-summary.json scripts/test3-load-pattern.js | tee results/test3-vu-output.txt
```
- **관찰(발췌)**: PASS — `p(95)=243.88ms`, `http_reqs rate=43.154389/s`
- **해석/의미**
  - VU 기반은 “동시 실행자 수”를 올리는 방식이라, 요청률은 시스템/스크립트(think time 포함)에 의해 **결과로서 결정**된다.

### 3-B) arrival-rate 기반 램프(ramping-arrival-rate)
```bash
MODE=rps THINK_TIME_S=0 k6 run --summary-export results/test3-rps-summary.json scripts/test3-load-pattern.js | tee results/test3-rps-output.txt
```
- **관찰(발췌)**: PASS — `p(95)=222.94ms`, `http_reqs rate=43.516803/s`
- **해석/의미**
  - arrival-rate 기반은 “목표 요청률”을 맞추려 하고, 부족하면 VU를 늘려 따라간다.  
  - 따라서 CAT에서 “RPS로 시나리오를 정의”하려면 이 모델이 더 직접적이다(달성 여부는 별도 확인 필요).

### 3-C) threshold 변경 실험(의도적 FAIL)
```bash
MODE=rps SLO_P95_MS=50 k6 run --summary-export results/test3-rps-fail-summary.json scripts/test3-load-pattern.js | tee results/test3-rps-fail-output.txt
```
- **관찰(발췌)**: FAIL (exit_code=99) — `p(95)<50` 위반, `p(95)=222.44ms`
- **해석**
  - threshold는 “참고 지표”가 아니라, 실행 결과를 **이진 판정(PASS/FAIL)**으로 만드는 장치다.
- **의미(CAT 관점)**
  - CAT 파이프라인에서 k6는 “부하 발생 + SLO 게이트”까지 **단독 수행**할 수 있다(최소 판정 엔진 역할 가능).

### 3-D) think time 영향(동일 VU 램프 + sleep)
```bash
MODE=vu THINK_TIME_S=0.2 k6 run --summary-export results/test3-vu-think-summary.json scripts/test3-load-pattern.js | tee results/test3-vu-think-output.txt
```
- **관찰(발췌)**: PASS — `p(95)=253.21ms`, `http_reqs rate=20.446724/s`
- **해석**
  - sleep이 늘면 VU당 iteration이 줄고, 결과적으로 `http_reqs rate`가 감소한다(“느려진 게 아니라 덜 보낸 것”).
- **의미(CAT 관점)**
  - 같은 VU 수라도 “사용자 사고시간”을 넣으면 트래픽이 줄 수 있으므로, **부하 목표가 VU인지 RPS인지**를 테스트 정의에 명시해야 한다.
