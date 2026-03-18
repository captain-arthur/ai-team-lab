import argparse
import json
from datetime import datetime, timezone


def iso_now():
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def get_metric(metrics, name):
    m = metrics.get(name)
    if m is None:
        return None
    return m


def main():
    p = argparse.ArgumentParser(description="Convert k6 summary-export JSON to CAT result JSON")
    p.add_argument("--k6-summary", required=True, help="Path to k6-summary.json (from --summary-export)")
    p.add_argument("--out", required=True, help="Output path for cat-result.json")
    p.add_argument("--test-name", required=True)
    p.add_argument("--scenario-type", required=True, help="CAT-level scenario type, e.g. http")
    p.add_argument("--target", required=True)
    p.add_argument("--exit-code", type=int, required=True)
    p.add_argument("--slo-latency-p95-ms", type=float, required=True)
    p.add_argument("--slo-error-rate", type=float, required=True)
    args = p.parse_args()

    with open(args.k6_summary, "r", encoding="utf-8") as f:
        summary = json.load(f)

    metrics = summary.get("metrics", {})

    http_req_duration = get_metric(metrics, "http_req_duration") or {}
    http_req_failed = get_metric(metrics, "http_req_failed") or {}
    http_reqs = get_metric(metrics, "http_reqs") or {}

    sli = {
        "latency_p95_ms": http_req_duration.get("p(95)"),
        "error_rate": http_req_failed.get("value"),
        "throughput_rps": http_reqs.get("rate"),
    }

    final_pass_fail = "PASS" if args.exit_code == 0 else "FAIL"

    result = {
        "test_name": args.test_name,
        "tool": "k6",
        "scenario_type": args.scenario_type,
        "target": args.target,
        "sli": sli,
        "slo": {
            "latency_p95_ms": args.slo_latency_p95_ms,
            "error_rate": args.slo_error_rate,
        },
        "final_pass_fail": final_pass_fail,
        "exit_code": args.exit_code,
        "timestamp": iso_now(),
        "artifacts": {
            "k6_summary_path": args.k6_summary,
        },
    }

    with open(args.out, "w", encoding="utf-8") as f:
        json.dump(result, f, ensure_ascii=False, indent=2)
        f.write("\n")


if __name__ == "__main__":
    main()

