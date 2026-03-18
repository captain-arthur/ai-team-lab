#!/usr/bin/env python3
import argparse
import json
import sys


def _read_json(path: str) -> dict:
    try:
        with open(path, "r", encoding="utf-8") as f:
            return json.load(f)
    except Exception as e:
        raise RuntimeError(f"json 읽기 실패: {path}: {e}") from e


def decide(cat_result: dict, monitoring_state: dict) -> dict:
    cat = cat_result.get("final_pass_fail")
    mon = monitoring_state.get("state")

    if cat not in ("PASS", "FAIL"):
        raise ValueError(f"cat.final_pass_fail 값이 필요합니다(PASS/FAIL). 현재: {cat!r}")
    if mon not in ("safe", "warning", "fail"):
        raise ValueError(f"monitoring.state 값이 필요합니다(safe/warning/fail). 현재: {mon!r}")

    # Rule 1 (FAIL 우선, monitoring과 무관)
    if cat == "FAIL":
        return {
            "system_state": "FAIL",
            "reason": "CAT FAIL",
            "source": {"cat": cat, "monitoring": mon},
        }

    # cat == "PASS" cases
    if mon == "fail" or mon == "warning":
        return {
            "system_state": "DEGRADED",
            "reason": f"CAT PASS and monitoring {mon}",
            "source": {"cat": cat, "monitoring": mon},
        }
    if mon == "safe":
        return {
            "system_state": "SAFE",
            "reason": "CAT PASS and monitoring safe",
            "source": {"cat": cat, "monitoring": mon},
        }

    # 위 검증에서(mon 값 제한) 도달 불가능
    raise RuntimeError("unreachable")


def main() -> int:
    p = argparse.ArgumentParser(
        description="CAT + Monitoring 통합 Decision Layer 최소 엔진"
    )
    p.add_argument("cat_result_json", help="CAT 결과 json 파일(예: cat-result.json)")
    p.add_argument("monitoring_state_json", help="Monitoring 상태 json 파일(예: monitoring-state.json)")
    p.add_argument("out_decision_json", help="출력 decision.json 경로")
    args = p.parse_args()

    cat_result = _read_json(args.cat_result_json)
    monitoring_state = _read_json(args.monitoring_state_json)
    decision = decide(cat_result, monitoring_state)

    with open(args.out_decision_json, "w", encoding="utf-8") as f:
        json.dump(decision, f, ensure_ascii=False, indent=2)
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        sys.exit(1)

