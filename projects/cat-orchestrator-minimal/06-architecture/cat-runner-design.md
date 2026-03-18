# Runner 설계(단순 CLI, 파일 기반)

## 1) 입력
- `job.yaml` (CAT Job Spec)

## 2) 처리 흐름(고정)
1. job.yaml 파싱
2. `tool` 선택
3. adapter 실행
   - adapter는 job spec을 입력으로 받고,
   - job 실행 후 `cat-result.json`을 `output.dir`에 생성
4. Runner는 `output.dir/cat-result.json` 존재 여부만 확인
5. Runner는 overall 파일(선택)을 생성할 수 있지만, 이번 설계에서는 필수 아님

## 3) 출력
- `output.dir/cat-result.json`
- adapter가 저장한 raw 결과(같은 output.dir 아래)

## 4) 실행 방식(단순 CLI)
```bash
cat run ./job.yaml
```

## 5) 의사코드(매우 단순)
```text
job = parse(job.yaml)
adapter = select(job.tool)
exit = adapter.run(job)
assert file_exists(job.output.dir + "/cat-result.json")
return cat-result.json (runner는 재판정하지 않음)
```

